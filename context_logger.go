package lumina

import (
    "context"
    "log"
)

type loggerContextKeyType struct{}

// A context key for connection-specific logger. The associated value will be of
// type *log.Logger.
var loggerContextKey = loggerContextKeyType{}

func GetLogger(ctx context.Context) *log.Logger {
    return ctx.Value(loggerContextKey).(*log.Logger)
}

func setLogger(ctx context.Context, logger *log.Logger) context.Context {
    return context.WithValue(ctx, loggerContextKey, logger)
}
