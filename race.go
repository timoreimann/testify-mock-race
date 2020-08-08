package race

import (
	"context"
	"fmt"
)

type Tx struct {
	ctx  context.Context
	done bool
}

func (tx *Tx) awaitDone() {
	go func() {
		<-tx.ctx.Done()
		tx.done = true
	}()
}

type Dependency interface {
	End(ctx context.Context, tx *Tx) error
	Fail(ctx context.Context, tx *Tx) error
}

type server struct {
	dep Dependency
}

func (d *server) Unsubscribe(ctx context.Context) error {
	txCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	tx := &Tx{ctx: txCtx}
	tx.awaitDone()

	if err := d.dep.End(ctx, tx); err != nil {
		return fmt.Errorf("End failed: %s", err)
	}
	if err := d.dep.Fail(ctx, tx); err != nil {
		return fmt.Errorf("Fail failed: %s", err)
	}
	return nil
}
