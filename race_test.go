package race

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var txType = mock.AnythingOfType("*race.Tx")

type EnderMock struct {
	mock.Mock
}

func (m *EnderMock) End(ctx context.Context, tx *Tx) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

func TestDo(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	myMock := &EnderMock{}
	myMock.On("End", ctx, txType).Return(nil)

	d := server{dep: myMock}
	err := d.Unsubscribe(ctx)
	assert.NoError(t, err)
	cancel()

	myMock.AssertExpectations(t)
}
