package storemock

import (
	"context"
	"database/sql"

	"github.com/genudine/saerro-go/types"
	"github.com/stretchr/testify/mock"
)

type MockVehicleStore struct {
	mock.Mock

	DB           *sql.DB
	RanMigration bool
}

func (m *MockVehicleStore) IsMigrated(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

func (m *MockVehicleStore) RunMigration(ctx context.Context, force bool) {
	m.Called(ctx, force)
}

func (m *MockVehicleStore) Insert(ctx context.Context, vehicle *types.Vehicle) error {
	args := m.Called(ctx, vehicle)
	return args.Error(0)
}

func (m *MockVehicleStore) GetOne(ctx context.Context, id string) (*types.Vehicle, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*types.Vehicle), args.Error(1)
}

func (m *MockVehicleStore) Prune(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return int64(args.Int(0)), args.Error(1)
}
