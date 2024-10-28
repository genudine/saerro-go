package storemock

import (
	"context"
	"database/sql"

	"github.com/genudine/saerro-go/types"
	"github.com/stretchr/testify/mock"
)

type MockPlayerStore struct {
	mock.Mock

	DB           *sql.DB
	RanMigration bool
}

func (m *MockPlayerStore) IsMigrated(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

func (m *MockPlayerStore) RunMigration(ctx context.Context, force bool) {
	m.Called(ctx, force)
}

func (m *MockPlayerStore) Insert(ctx context.Context, player *types.Player) error {
	args := m.Called(ctx, player)
	return args.Error(0)
}

func (m *MockPlayerStore) GetOne(ctx context.Context, id string) (*types.Player, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*types.Player), args.Error(1)
}

func (m *MockPlayerStore) Prune(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return int64(args.Int(0)), args.Error(1)
}
