package ctxutil

import (
	"context"
	"fmt"
)

func CreateContext(ctx context.Context, key ContextKey, value any) context.Context {
	return context.WithValue(ctx, key, value)
}

func GetContext(ctx context.Context, key ContextKey) (any, error) {
	val := ctx.Value(key)
	if val == nil {
		return nil, fmt.Errorf("Context key not found")
	}

	return ctx.Value(key), nil
}
