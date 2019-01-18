package lumina

import "context"

type Handler interface {
    Interpreter

    // Check whether the incoming request is accepted, before unmarshalling.
    AcceptRequest(t PacketType) bool

    // ServeRequest will be called with unmarshalled and validated Request.
    //
    // For proxy senario, if any fields of the request was changed and the
    // request contains a valid cache, you can assume that req implements
    // CacheWriter interface and call SetCache(nil) on it. The packet fields
    // will be re-validated and then marshalled by Client.
    ServeRequest(ctx context.Context, req Request) (rsp Packet, err error)
}
