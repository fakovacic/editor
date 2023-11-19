package app

import "context"

const (
	RequestID Key = "reqID"
)

type Key string

func (k Key) String() string {
	return string(k)
}

func GetStringValue(ctx context.Context, key Key) string {
	ctxValue := ctx.Value(key)

	if ctxValue != nil {
		val, ok := ctxValue.(string)
		if ok {
			return val
		}
	}

	return ""
}
