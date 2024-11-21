package xcontext

import (
	"context"
	"github.com/google/uuid"
	"strings"
)

const (
	traceIdKey     = "trace-id"
	defaultTraceId = "0000000000000000"
)

func GetTraceId(ctx context.Context) string {
	if traceId, ok := ctx.Value(traceIdKey).(string); ok {
		return traceId
	}
	return defaultTraceId
}

func SetTraceId(ctx context.Context) context.Context {
	if _, ok := ctx.Value(traceIdKey).(string); !ok {
		return context.WithValue(ctx, traceIdKey, strings.ReplaceAll(uuid.New().String(), "-", "")[0:16])
	}
	return ctx
}

func TraceIdKey() string {
	return traceIdKey
}
