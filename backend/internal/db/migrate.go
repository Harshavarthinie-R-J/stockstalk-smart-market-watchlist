package db

import (
	"context"
	"database/sql"
	"os"
)

// Schema setup.
func Migrate(ctx context.Context, conn *sql.DB, path string) error {
	contents, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_, err = conn.ExecContext(ctx, string(contents))
	return err
}
