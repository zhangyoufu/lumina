package lumina

import (
    "context"
    "net"
)

type connContextKeyType struct{}

// A context key for connection-specific logger. The associated value will be of
// type net.Conn.
var connContextKey = connContextKeyType{}

func GetConn(ctx context.Context) net.Conn {
    return ctx.Value(connContextKey).(net.Conn)
}

func setConn(ctx context.Context, conn net.Conn) context.Context {
    return context.WithValue(ctx, connContextKey, conn)
}
