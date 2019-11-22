package context_util

import (
	"context"
	"errors"
)

const (
	ContextKeyUserID ContextKey = iota
)

type ContextKey int

func WithUserID(c context.Context, userID int64) context.Context {
	return context.WithValue(c, ContextKeyUserID, userID)
}

func GetUserID(c context.Context) (int64, error) {
	val := c.Value(ContextKeyUserID)
	if val == nil {
		return 0, errors.New("User ID not injected yet")
	}
	return val.(int64), nil
}
