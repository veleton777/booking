package transaction

import (
	"context"
)

const TxKey = "transaction"

//go:generate mockery --name Tx
type Tx interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
