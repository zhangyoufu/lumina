package main

import (
    "context"
    "github.com/zhangyoufu/lumina"
)

type upstreamContextKeyType struct{}

// A context key for upstream connection. The associated value will be of type
// *lumina.ClientSession.
var upstreamContextKey = upstreamContextKeyType{}

// Return *lumina.ClientSession extracted from a given context.Context.
func getUpstream(ctx context.Context) *lumina.ClientSession {
    return ctx.Value(upstreamContextKey).(*lumina.ClientSession)
}

// Create a lumina.ClientSession instance as upstream for each incoming connection.
func setUpstream(ctx context.Context, session *lumina.ClientSession) context.Context {
    return context.WithValue(ctx, upstreamContextKey, session)
}
