package optimizer

import (
	"github.com/stretchr/testify/mock"
)

var _ Optimizer = (*MockOptimizer)(nil)

type MockOptimizer struct {
	mock.Mock
}

func NewMockOptimizer() *MockOptimizer {
	return &MockOptimizer{}
}

func (r *MockOptimizer) Execute(filenames string) (*Result, error) {
	args := r.Called(filenames)
	return args.Get(0).(*Result), args.Error(1)
}
