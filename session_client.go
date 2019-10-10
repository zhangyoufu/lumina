package lumina

import (
	"bytes"
	"context"
	"encoding/hex"
	"github.com/palantir/stacktrace"
	"log"
	"net"
)

var _ = hex.EncodeToString

type ClientSession struct {
	conn        net.Conn
	logger      *log.Logger
	interpreter Interpreter
}

func (c *Client) Dial(ctx context.Context, logger *log.Logger, interpreter Interpreter) (s *ClientSession, err error) {
	// FIXME: Dialer does not support context.Context
	conn, err := c.getDialer().Dial()
	if err != nil {
		err = stacktrace.Propagate(err, "dial failed")
		return
	}

	if logger == nil {
		logger = newTaggedLogger()
	}

	_s := &ClientSession{
		conn:        conn,
		logger:      logger,
		interpreter: interpreter,
	}

	rsp, err := _s.Request(ctx, newHeloPacket(c.LicenseKey, c.LicenseId))
	if err != nil {
		err = stacktrace.Propagate(err, "hello request failed")
		conn.Close()
		return
	}
	if _, ok := rsp.(*RpcOkPacket); !ok {
		err = stacktrace.NewError("hello response: %#v", rsp)
		conn.Close()
		return
	}

	// The underlying net.Conn has already called runtime.SetFinalizer.
	s = _s
	return
}

func (s *ClientSession) Close() error {
	return s.conn.Close()
}

// The packet cache will be used if available.
func (s *ClientSession) Request(ctx context.Context, req Request) (rsp Packet, err error) {
	reqType := req.getType()
	s.logger.Print("client send request type: ", reqType)
	// s.logger.Printf("client send request data: %+v", req)
	var reqRaw RawPacket
	if cr, _ := req.(cacheReader); cr != nil {
		reqRaw = cr.getCache()
	}
	if reqRaw == nil {
		if err = req.validateFields(); err != nil {
			err = stacktrace.Propagate(err, "validation failed for request of type %v", reqType)
			return
		}
		reqPayload := &bytes.Buffer{}
		if err = req.WriteTo(reqPayload); err != nil {
			err = stacktrace.Propagate(err, "unable to marshal request of type %v", reqType)
			return
		}
		reqRaw, err = NewRawPacket(reqType, reqPayload.Bytes())
		if err != nil {
			err = stacktrace.Propagate(err, "unable to compose raw packet for request of type %v", reqType)
			return
		}
	}
	// s.logger.Print("client send request rawdata: ", hex.EncodeToString(reqRaw.GetPayload()))
	if err = reqRaw.WriteTo(s.conn); err != nil {
		err = stacktrace.Propagate(err, "unable to write raw packet to server")
		s.Close()
		return
	}

	var rspRaw RawPacket
	err = rspRaw.ReadFrom(s.conn)
	if err != nil {
		err = stacktrace.Propagate(err, "unable to read raw packet from server")
		s.Close()
		return
	}
	rspType := rspRaw.GetType()
	s.logger.Print("client recv response type: ", rspType)
	// s.logger.Print("client recv response rawdata: ", hex.EncodeToString(rspRaw.GetPayload()))
	rsp = rpcInterpreter.GetPacketOfType(rspType)
	if rsp == nil {
		rsp = s.interpreter.GetPacketOfType(rspType)
	}
	if rsp == nil {
		err = stacktrace.NewError("response of type %v is not supported by interpreter", rspType)
		return
	}
	r := bytes.NewReader(rspRaw.GetPayload())
	err = rsp.ReadFrom(r)
	if err == nil && r.Len() > 0 {
		err = errTrailingData
	}
	if err != nil {
		s.logger.Print(stacktrace.Propagate(err, "unable to unmarshal response of type %v", rspType))
		return
	}
	// s.logger.Printf("client recv response data: %+v", rsp)
	if err = rsp.validateFields(); err != nil {
		err = stacktrace.Propagate(err, "validation failed for response of type %v", rspType)
		return
	}
	if cw, _ := rsp.(cacheWriter); cw != nil {
		cw.setCache(rspRaw)
	}

	return
}
