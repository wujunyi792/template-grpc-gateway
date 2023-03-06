package contextx

import (
	"context"
	"pinnacle-primary-be/core/mappingx"
)

const contextTagKey = "ctx"

var unmarshaler = mappingx.NewUnmarshaler(contextTagKey)

type contextValuer struct {
	context.Context
}

func (cv contextValuer) Value(key string) (any, bool) {
	v := cv.Context.Value(key)
	return v, v != nil
}

// For unmarshals ctx into v.
func For(ctx context.Context, v any) error {
	return unmarshaler.UnmarshalValuer(contextValuer{
		Context: ctx,
	}, v)
}
