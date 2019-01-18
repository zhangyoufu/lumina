package main

import (
    "context"
    "log"
    "net"
    "github.com/zhangyoufu/lumina"
)

type Handler struct {}

func (*Handler) GetPacketOfType(t lumina.PacketType) lumina.Packet {
    switch t {
    case lumina.PKT_PULL_MD:
        return &lumina.PullMdPacket{}
    case lumina.PKT_PUSH_MD:
        return &lumina.PushMdPacket{}
    default:
        return nil
    }
}

func (*Handler) AcceptRequest(lumina.PacketType) bool {
    return true
}

func (*Handler) ServeRequest(ctx context.Context, request lumina.Request) (response lumina.Packet, err error) {
    _ = ctx
    switch request.(type) {
    case *lumina.PullMdPacket:
        req := request.(*lumina.PullMdPacket)
        rsp := &lumina.PullMdResultPacket{}
        rsp.Codes = make([]lumina.OpResult, len(req.PatternIds))
        response = rsp
    case *lumina.PushMdPacket:
        req := request.(*lumina.PushMdPacket)
        logger := lumina.GetLogger(ctx)
        logger.Print("idb: ", req.Idb)
        logger.Print("input: ", req.Input)
        logger.Print("hostname: ", req.Hostname)
        rsp := &lumina.PushMdResultPacket{}
        rsp.Codes = make([]lumina.OpResult, len(req.Contents))
        response = rsp
    default:
        log.Fatal("unable to serve")
    }
    return
}

func main() {
    ln, err := net.ListenTCP("tcp", &net.TCPAddr{Port: 65432})
    if err != nil {
        log.Fatal("unable to listen: ", err)
    }
    srv := lumina.Server{ Handler: &Handler{} }
    log.Print("server is listening on ", ln.Addr())
    srv.Serve(ln)
}
