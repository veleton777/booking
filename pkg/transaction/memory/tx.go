package memory

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/veleton777/booking_api/pkg/transaction"
)

var ErrTransactionNotFoundInCtx = errors.New("transaction not found in context")

type TxClient struct {
	mu *sync.Mutex
}

func NewTxClient() *TxClient {
	return &TxClient{mu: &sync.Mutex{}}
}

func (t *TxClient) Begin(ctx context.Context) (context.Context, error) {
	ctx = context.WithValue(ctx, transaction.TxKey, uuid.New())

	t.mu.Lock()

	return ctx, nil
}

func (t *TxClient) Commit(ctx context.Context) error {
	err := t.txFromCtx(ctx)
	if err != nil {
		return errors.Wrap(err, "get tx from context")
	}

	t.mu.Unlock()

	return nil
}

func (t *TxClient) Rollback(ctx context.Context) error {
	err := t.txFromCtx(ctx)
	if err != nil {
		return errors.Wrap(err, "get tx from context")
	}

	t.mu.TryLock()
	t.mu.Unlock()

	// todo logic for rollback

	return nil
}

func (t *TxClient) TxFunc(ctx context.Context, f func(context.Context) error) error {
	ctx, err := t.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}

	err = f(ctx)
	if err != nil {
		defer t.Rollback(ctx)

		return errors.Wrap(err, "call func")
	}

	if err = t.Commit(ctx); err != nil {
		return errors.Wrap(err, "commit transaction")
	}

	return nil
}

func (t *TxClient) txFromCtx(ctx context.Context) error {
	v := ctx.Value(transaction.TxKey)

	_, ok := v.(uuid.UUID)
	if !ok {
		return ErrTransactionNotFoundInCtx
	}

	return nil
}
