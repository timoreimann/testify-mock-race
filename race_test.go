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

func (m *EnderMock) Fail(ctx context.Context, tx *Tx) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

func TestDo(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	myMock := &EnderMock{}
	myMock.On("End", ctx, txType).Return(nil)
	// Commenting in this one makes the race go away :(
	// myMock.On("Fail", ctx, txType).Return(errors.New("failing hard"))

	d := server{dep: myMock}
	err := d.Unsubscribe(ctx)
	assert.Error(t, err)
	cancel()

	myMock.AssertExpectations(t)
}
