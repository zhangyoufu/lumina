package lumina

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"github.com/palantir/stacktrace"
	"log"
	"net"
)

var _ = hex.EncodeToString

type HeloCallback func(context.Context) (context.Context, error)

type Server struct {
	Handler Handler

	// Called once after PKT_HELO is received. Protocol processing is deferred
	// until the callback returns.
	OnHELO HeloCallback
}

func (srv *Server) Serve(ln net.Listener) error {
	if srv.Handler == nil {
		return stacktrace.NewError("server handler not set")
	}

	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			return stacktrace.Propagate(err, "unable to accept")
		}
		go srv.serveConn(conn)
	}
}

// Shutdown gracefully shuts down the server without interrupting any active
// connections. Shutdown works by first closing all open listeners, then closing
// all idle connections, and then waiting indefinitely for connections to return
// to idle and then shut down. If the provided context expires before the shut-
// down is complete, Shutdown returns the context's error, otherwise it returns
// any error returned from closing the Server's underlying Listener(s).
//
// When Shutdown is called, Serve immediately return ErrServerClosed. Make sure
// the program doesn't exit and waits instead for Shutdown to return.
//
// Once Shutdown has been called on a server, it may not be reused; future calls
// to methods such as Serve will return ErrServerClosed.
// func (s *Server) Shutdown(ctx context.Context) error {
//     return nil
// }

func (srv *Server) serveConn(conn net.Conn) {
	defer conn.Close()

	logger := newTaggedLogger()
	logger.Print("incoming connection from ", conn.RemoteAddr())

	ctx := context.TODO()
	ctx = setConn(ctx, conn)
	ctx = setLogger(ctx, logger)

	s := ServerSession{
		conn:   conn,
		logger: logger,
		ctx:    ctx,
	}
	s.serve(srv.Handler, srv.OnHELO)
}

type ServerSession struct {
	conn    net.Conn
	logger  *log.Logger
	ctx     context.Context
	handler Handler
}

var errInternalServerError = errors.New("internal server error")

func (s *ServerSession) serveOne(handler Handler) error {
	var (
		err error
		req Request
		rsp Packet
	)
	req, err = s.recvRequest(handler)
	switch err {
	case nil:
		// TODO: further refine context
		rsp, err = handler.ServeRequest(s.ctx, req)
		if err != nil {
			err = stacktrace.Propagate(err, "error while serving request")
		} else {
			reqType := req.getType()
			rspType := rsp.getType()
			if rspType != req.getResponseType() && rspType != PKT_RPC_FAIL {
				err = stacktrace.NewError("%s response is not allowed for %s request", rspType, reqType)
			}
		}
	case errInternalServerError: // non-critical
		break
	default: // critical
		return stacktrace.Propagate(err, "unrecoverable error while receiving request")
	}

	if err != nil {
		s.logger.Print(err)
		rsp = &RpcFailPacket{
			Result: -1,
			Error:  "internal server error",
		}
	}

	if err = s.sendResponse(rsp); err != nil {
		return stacktrace.Propagate(err, "unrecoverable error while sending response")
	}

	return nil
}

func (s *ServerSession) serve(handler Handler, onHelo HeloCallback) {
	var err error
	var newCtx context.Context

	if err = s.serveOne(heloHandler); err != nil {
		goto _ret
	}

	if onHelo != nil {
		if newCtx, err = onHelo(s.ctx); err != nil {
			goto _ret
		}
		s.ctx = newCtx
	}

	for {
		if err = s.serveOne(handler); err != nil {
			goto _ret
		}
	}

_ret:
	s.logger.Print(err)
	return
}

func (s *ServerSession) recvRequest(h Handler) (req Request, err error) {
	var reqRaw RawPacket
	err = reqRaw.readFrom(s.conn)
	if err != nil {
		err = stacktrace.Propagate(err, "unable to read raw packet from client")
		return
	}
	// TODO: return a SERVFAIL instead of abort connection
	reqType := reqRaw.GetType()
	s.logger.Print("server recv request type: ", reqType)
	// logger.Print("server recv request rawdata: ", hex.EncodeToString(reqRaw.GetPayload()))
	if !h.AcceptRequest(reqType) {
		s.logger.Printf("request of type %v is not accepted by server", reqType)
		err = errInternalServerError
		return
	}
	req, _ = h.GetPacketOfType(reqType).(Request)
	if req == nil {
		s.logger.Printf("request of type %v is not supported by server", reqType)
		err = errInternalServerError
		return
	}
	r := bytes.NewReader(reqRaw.GetPayload())
	err = req.readFrom(r)
	if err == nil && r.Len() > 0 {
		err = errTrailingData
	}
	if err != nil {
		s.logger.Print(stacktrace.Propagate(err, "unable to unmarshal request of type %v", reqType))
		err = errInternalServerError
		return
	}
	// logger.Printf("server recv request data: %+v", req)
	if err = req.validateFields(); err != nil {
		s.logger.Print(stacktrace.Propagate(err, "validation failed for request of type %v", reqType))
		err = errInternalServerError
		return
	}
	if cw, _ := req.(cacheWriter); cw != nil {
		cw.setCache(reqRaw)
	}
	return
}

// TODO: standalone code path for internal server error
func (s *ServerSession) sendResponse(rsp Packet) (err error) {
	rspType := rsp.getType()
	s.logger.Print("server send response type: ", rspType)
	// logger.Printf("server send response data: %+v", rsp)
	var rspRaw RawPacket
	if cr, _ := rsp.(cacheReader); cr != nil {
		rspRaw = cr.getCache()
	}
	if rspRaw == nil {
		if err = rsp.validateFields(); err != nil {
			s.logger.Print(stacktrace.Propagate(err, "validation failed for response of type %v", rspType))
			return
		}
		rspPayload := &bytes.Buffer{}
		if err = rsp.writeTo(rspPayload); err != nil {
			s.logger.Print(stacktrace.Propagate(err, "unable to marshal response of type %v", rspType))
			return
		}
		rspRaw, err = NewRawPacket(rspType, rspPayload.Bytes())
		if err != nil {
			s.logger.Print(stacktrace.Propagate(err, "unable to compose raw packet for response of type %v", rspType))
			return
		}
	}
	// logger.Print("server send response rawdata: ", hex.EncodeToString(rspRaw.GetPayload()))
	if err = rspRaw.writeTo(s.conn); err != nil {
		s.logger.Print(stacktrace.Propagate(err, "unable to write raw packet to client"))
		return
	}
	return
}
