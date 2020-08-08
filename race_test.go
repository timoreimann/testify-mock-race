package race

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	myMock := NewMockDependency(ctrl)

	myMock.EXPECT().End(gomock.Any(), gomock.Any()).Return(nil)

	d := server{dep: myMock}
	err := d.Unsubscribe(ctx)
	assert.NoError(t, err)
	cancel()
}
