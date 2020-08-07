package race

import (
	"context"
	"fmt"
)

type Tx struct {
	val  int
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
	res int
}

func (d *server) Unsubscribe(ctx context.Context) error {
	txCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	tx := &Tx{val: 22, ctx: txCtx}
	tx.awaitDone()

	d.res = tx.val

	if err := d.dep.End(ctx, tx); err != nil {
		return fmt.Errorf("End failed: %s", err)
	}
	if err := d.dep.Fail(ctx, tx); err != nil {
		return fmt.Errorf("Fail failed: %s", err)
	}
	return nil
}
