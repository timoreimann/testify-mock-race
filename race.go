package race

import (
	"context"
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
}

type server struct {
	dep Dependency
}

func (d *server) Unsubscribe(ctx context.Context) error {
	txCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	tx := &Tx{ctx: txCtx}
	tx.awaitDone()

	return d.dep.End(ctx, tx)
}
