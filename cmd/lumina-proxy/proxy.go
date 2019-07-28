package main

import (
    "context"
    "github.com/palantir/stacktrace"
    "github.com/zhangyoufu/lumina"
)

// 'W' for Windows
// 'M' for macOS
// 'L' for Linux
var licOS byte

type Proxy struct {
    lumina.Server
}

func NewProxy(licKey lumina.LicenseKey, licId lumina.LicenseId) (proxy *Proxy) {
    client := &lumina.Client{
        LicenseKey: licKey,
        LicenseId: licId,
    }
    handler := &proxyHandler{}
    proxy = &Proxy{}
    proxy.Handler = handler
    proxy.OnHELO = func(ctx context.Context) (newctx context.Context, err error) {
        session, err := client.Dial(ctx, lumina.GetLogger(ctx), handler)
        if err != nil {
            if err != context.Canceled {
                err = stacktrace.Propagate(err, "unable to create upstream session")
            }
            return
        }
        newctx = setUpstream(ctx, session)
        return
    }
    return
}

type proxyHandler struct {}

// Currently, we only allow pulling metadata from upstream lumina server to
// avoid abuse of our proxy.
func (*proxyHandler) AcceptRequest(t lumina.PacketType) bool {
    switch t {
    case lumina.PKT_PULL_MD:
        return true
    default:
        return false
    }
}

func (*proxyHandler) GetPacketOfType(t lumina.PacketType) lumina.Packet {
    switch t {
    case lumina.PKT_PULL_MD:
        return &lumina.PullMdPacket{}
    case lumina.PKT_PULL_MD_RESULT:
        return &lumina.PullMdResultPacket{}
    default:
        return nil
    }
}

// Pump between client and upstream server. (half-duplex)
func (*proxyHandler) ServeRequest(ctx context.Context, req lumina.Request) (rsp lumina.Packet, err error) {
    rsp, err = getUpstream(ctx).Request(ctx, req)
    // if err != nil {
    //     lumina.GetConn(ctx).Close()
    // }
    return
}
