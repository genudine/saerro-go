package store_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/genudine/saerro-go/util/testutil"
)

func mkEnv(t *testing.T) (context.Context, *sql.DB) {
	t.Helper()

	db := testutil.GetTestDB(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	t.Cleanup(func() {
		cancel()
	})

	return ctx, db
}
