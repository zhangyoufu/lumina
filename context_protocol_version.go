package lumina

import "context"

type protocolVersionContextKeyType struct{}

// A context key for lumina protocol version.
var protocolVersionContextKey = protocolVersionContextKeyType{}

// Return lumina protocol version extracted from a given context.Context.
func GetProtocolVersion(ctx context.Context) int32 {
	return ctx.Value(protocolVersionContextKey).(int32)
}

// Return a new context with given lumina protocol version.
func setProtocolVersion(ctx context.Context, version int32) context.Context {
	return context.WithValue(ctx, protocolVersionContextKey, version)
}
