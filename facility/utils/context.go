package utils

import (
	"context"
	"database/sql"
)

const (
	contextKeyTx contextKey = iota
)

type contextKey int

func ContextWithTx(tx *sql.Tx) context.Context {
	return context.WithValue(context.Background(), contextKeyTx, tx)
}

func GetTx(ctx context.Context) *sql.Tx {
	val := ctx.Value(contextKeyTx)
	if tx, ok := val.(*sql.Tx); ok {
		return tx
	}
	return nil
}
