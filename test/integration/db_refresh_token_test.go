//go:build integration
// +build integration

package integration

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/ArtemVoronov/indefinite-studies-api/internal/db"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db/queries"
	"github.com/stretchr/testify/assert"
)

func TestDBRefreshTokenGet(t *testing.T) {
	t.Run("NotFoundCase", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			_, err := queries.GetRefreshTokenByToken(tx, ctx, TEST_REFRESH_TOKEN_1)

			assert.Equal(t, sql.ErrNoRows, err)
			return err
		})()
	})))
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		expected := utils.entityGenerators.GenerateRefreshToken(1)
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.CreateRefreshToken(tx, ctx, expected.UserId, expected.Token, expected.ExpireAt)

			assert.Nil(t, err)

			return err
		})()
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			actual, err := queries.GetRefreshTokenByToken(tx, ctx, expected.Token)

			utils.asserts.AssertEqualRefreshTokens(t, expected, actual)

			return err
		})()
	})))
	t.Run("TimeoutError", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at loading refresh token '%s' from db, case after QueryRow.Scan: %s", TEST_REFRESH_TOKEN_1, "context deadline exceeded")
			_, err := tx.ExecContext(ctx, "SELECT pg_sleep(10)")
			_, err = queries.GetRefreshTokenByToken(tx, ctx, TEST_REFRESH_TOKEN_1)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
	t.Run("ContextCancelled", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at loading refresh token '%s' from db, case after QueryRow.Scan: %s", TEST_REFRESH_TOKEN_1, "context canceled")
			cancel()
			_, err := queries.GetRefreshTokenByToken(tx, ctx, TEST_REFRESH_TOKEN_1)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
}

func TestDBRefreshTokenCreate(t *testing.T) {
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.CreateRefreshToken(tx, ctx, TEST_REFRESH_TOKEN_USER_ID_1, TEST_REFRESH_TOKEN_1, TEST_REFRESH_TOKEN_EXPIRE_AT_1)

			assert.Nil(t, err)
			return err
		})()
	})))
	t.Run("TimeoutError", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at creating refresh token, case after preparing statement: %s", "context deadline exceeded")
			_, err := tx.ExecContext(ctx, "SELECT pg_sleep(10)")
			err = queries.CreateRefreshToken(tx, ctx, TEST_REFRESH_TOKEN_USER_ID_1, TEST_REFRESH_TOKEN_1, TEST_REFRESH_TOKEN_EXPIRE_AT_1)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
	t.Run("ContextCancelled", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at creating refresh token, case after preparing statement: %s", "context canceled")
			cancel()
			err := queries.CreateRefreshToken(tx, ctx, TEST_REFRESH_TOKEN_USER_ID_1, TEST_REFRESH_TOKEN_1, TEST_REFRESH_TOKEN_EXPIRE_AT_1)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
}

func TestDBRefreshTokenUpdate(t *testing.T) {
	t.Run("NotFoundCase", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.UpdateRefreshToken(tx, ctx, TEST_REFRESH_TOKEN_USER_ID_2, TEST_REFRESH_TOKEN_2, TEST_REFRESH_TOKEN_EXPIRE_AT_2)

			assert.Equal(t, sql.ErrNoRows, err)
			return err
		})()
	})))
	t.Run("DeletedCase", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.CreateRefreshToken(tx, ctx, TEST_REFRESH_TOKEN_USER_ID_1, TEST_REFRESH_TOKEN_1, TEST_REFRESH_TOKEN_EXPIRE_AT_1)

			assert.Nil(t, err)

			return err
		})()
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.DeleteRefreshToken(tx, ctx, TEST_REFRESH_TOKEN_USER_ID_1)

			assert.Nil(t, err)

			return err
		})()
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.UpdateRefreshToken(tx, ctx, TEST_REFRESH_TOKEN_USER_ID_1, TEST_REFRESH_TOKEN_2, TEST_REFRESH_TOKEN_EXPIRE_AT_2)

			assert.Equal(t, sql.ErrNoRows, err)
			return err
		})()
	})))
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.CreateRefreshToken(tx, ctx, TEST_REFRESH_TOKEN_USER_ID_1, TEST_REFRESH_TOKEN_1, TEST_REFRESH_TOKEN_EXPIRE_AT_1)

			assert.Nil(t, err)

			return err
		})()
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.UpdateRefreshToken(tx, ctx, TEST_REFRESH_TOKEN_USER_ID_1, TEST_REFRESH_TOKEN_2, TEST_REFRESH_TOKEN_EXPIRE_AT_2)

			assert.Nil(t, err)

			return err
		})()
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			actual, err := queries.GetRefreshTokenByUserId(tx, ctx, TEST_REFRESH_TOKEN_USER_ID_1)

			assert.Equal(t, TEST_REFRESH_TOKEN_USER_ID_1, actual.UserId)
			assert.Equal(t, TEST_REFRESH_TOKEN_2, actual.Token)
			return err
		})()
	})))
	t.Run("TimeoutError", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at updating refresh token, case after preparing statement: %s", "context deadline exceeded")
			_, err := tx.ExecContext(ctx, "SELECT pg_sleep(10)")
			err = queries.UpdateRefreshToken(tx, ctx, TEST_REFRESH_TOKEN_USER_ID_1, TEST_REFRESH_TOKEN_2, TEST_REFRESH_TOKEN_EXPIRE_AT_2)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
	t.Run("ContextCancelled", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at updating refresh token, case after preparing statement: %s", "context canceled")
			cancel()
			err := queries.UpdateRefreshToken(tx, ctx, TEST_REFRESH_TOKEN_USER_ID_1, TEST_REFRESH_TOKEN_2, TEST_REFRESH_TOKEN_EXPIRE_AT_2)
			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
}

func TestDBRefreshTokenDelete(t *testing.T) {
	t.Run("NotFoundCase", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.DeleteRefreshToken(tx, ctx, TEST_REFRESH_TOKEN_USER_ID_1)

			assert.Equal(t, sql.ErrNoRows, err)
			return err
		})()
	})))
	t.Run("AlreadyDeletedCase", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.CreateRefreshToken(tx, ctx, TEST_REFRESH_TOKEN_USER_ID_1, TEST_REFRESH_TOKEN_1, TEST_REFRESH_TOKEN_EXPIRE_AT_1)

			assert.Nil(t, err)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.DeleteRefreshToken(tx, ctx, TEST_REFRESH_TOKEN_USER_ID_1)

			assert.Nil(t, err)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.DeleteRefreshToken(tx, ctx, TEST_REFRESH_TOKEN_USER_ID_1)

			assert.Equal(t, sql.ErrNoRows, err)
			return err
		})()
	})))
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.CreateRefreshToken(tx, ctx, TEST_REFRESH_TOKEN_USER_ID_1, TEST_REFRESH_TOKEN_1, TEST_REFRESH_TOKEN_EXPIRE_AT_1)

			return err
		})()
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.DeleteRefreshToken(tx, ctx, TEST_REFRESH_TOKEN_USER_ID_1)

			assert.Nil(t, err)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			_, err := queries.GetRefreshTokenByToken(tx, ctx, TEST_REFRESH_TOKEN_1)

			assert.Equal(t, sql.ErrNoRows, err)
			return err
		})()
	})))
	t.Run("TimeoutError", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at deleting refresh token, case after preparing statement: %s", "context deadline exceeded")
			_, err := tx.ExecContext(ctx, "SELECT pg_sleep(10)")
			err = queries.DeleteRefreshToken(tx, ctx, TEST_REFRESH_TOKEN_USER_ID_1)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
	t.Run("ContextCancelled", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at deleting refresh token, case after preparing statement: %s", "context canceled")
			cancel()
			err := queries.DeleteRefreshToken(tx, ctx, TEST_REFRESH_TOKEN_USER_ID_1)
			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
}
