package lumina

import (
    "bytes"
    "context"
    "encoding/hex"
    "errors"
    "log"
    "net"
    "github.com/palantir/stacktrace"
)

var _ = hex.EncodeToString

type Server struct {
    Handler     Handler

    // A callback for new connections. Protocol processing is deferred until the
    // callback returns.
    OnAccept    func(context.Context) (context.Context, error)
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

    logger := newConnLogger(conn)
    logger.Print("incoming connection from ", conn.RemoteAddr())

    ctx := context.TODO()
    ctx = setConn(ctx, conn)
    ctx = setLogger(ctx, logger)
    if srv.OnAccept != nil {
        var err error
        ctx, err = srv.OnAccept(ctx)
        if err != nil {
            logger.Print(stacktrace.Propagate(err, "error in OnAccept callback"))
            return
        }
    }

    s := ServerSession{
        conn: conn,
        logger: logger,
        ctx: ctx,
    }
    s.serve(srv.Handler)
}

type ServerSession struct {
    conn    net.Conn
    logger  *log.Logger
    ctx     context.Context
    handler Handler
}

var errInternalServerError = errors.New("internal server error")

func (s *ServerSession) serve(handler Handler) {
    s.handler = heloHandler
    for {
        var rsp Packet
        req, err := s.recvRequest()
        switch err {
        case nil:
            // TODO: further refine context
            rsp, err = s.handler.ServeRequest(s.ctx, req)
            if err != nil {
                err = stacktrace.Propagate(err, "error while serving request")
                s.logger.Print(err)
                rsp = nil
            } else {
                reqType := req.getType()
                rspType := rsp.getType()
                if rspType != req.getResponseType() && rspType != PKT_RPC_FAIL {
                    s.logger.Printf("%s response is not allowed for %s request", rspType, reqType)
                    return
                }
            }
        case errInternalServerError:
            // break
        default:
            stacktrace.Propagate(err, "unrecoverable error while receiving request")
            s.logger.Print(err)
            return
        }

        if rsp == nil {
            rsp = &RpcFailPacket{
                Result: -1,
                Error: "internal server error",
            }
        }

        if err = s.sendResponse(rsp); err != nil {
            err = stacktrace.Propagate(err, "unrecoverable error while sending request")
            s.logger.Print(err)
            return
        }

        s.handler = handler
    }
}

func (s *ServerSession) recvRequest() (req Request, err error) {
    var reqRaw RawPacket
    err = reqRaw.ReadFrom(s.conn)
    if err != nil {
        s.logger.Print(stacktrace.Propagate(err, "unable to read raw packet from client"))
        return
    }
    // TODO: return a SERVFAIL instead of abort connection
    reqType := reqRaw.GetType()
    s.logger.Print("server recv request type: ", reqType)
    // logger.Print("server recv request rawdata: ", hex.EncodeToString(reqRaw.GetPayload()))
    if !s.handler.AcceptRequest(reqType) {
        s.logger.Printf("request of type %v is not accepted by server", reqType)
        err = errInternalServerError
        return
    }
    req, _ = s.handler.GetPacketOfType(reqType).(Request)
    if req == nil {
        s.logger.Printf("request of type %v is not supported by server", reqType)
        err = errInternalServerError
        return
    }
    r := bytes.NewReader(reqRaw.GetPayload())
    err = req.ReadFrom(r)
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
        if err = rsp.WriteTo(rspPayload); err != nil {
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
    if err = rspRaw.WriteTo(s.conn); err != nil {
        s.logger.Print(stacktrace.Propagate(err, "unable to write raw packet to client"))
        return
    }
    return
}
