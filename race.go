package race

import (
	"context"
	"fmt"
)

type Tx struct {
	val int
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
	tx := &Tx{23}
	go func() {
		<-ctx.Done()
		d.res = tx.val
	}()

	if err := d.dep.End(ctx, tx); err != nil {
		return fmt.Errorf("End failed: %s", err)
	}
	if err := d.dep.Fail(ctx, tx); err != nil {
		return fmt.Errorf("Fail failed: %s", err)
	}
	return nil
}
