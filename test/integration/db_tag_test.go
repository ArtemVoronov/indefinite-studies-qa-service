//go:build integration
// +build integration

package integration

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/ArtemVoronov/indefinite-studies-api/internal/db"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db/entities"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db/queries"
	"github.com/stretchr/testify/assert"
)

func TestDBTagGet(t *testing.T) {
	t.Run("NotFoundCase", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			_, err := queries.GetTag(tx, ctx, 1)

			assert.Equal(t, sql.ErrNoRows, err)
			return err
		})()
	})))
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		expected := utils.entityGenerators.GenerateTag(1)
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			tagId, err := queries.CreateTag(tx, ctx, expected.Name, expected.State)
			assert.Nil(t, err)
			assert.Equal(t, tagId, expected.Id)
			return err
		})()
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			actual, err := queries.GetTag(tx, ctx, expected.Id)
			utils.asserts.AssertEqualTags(t, expected, actual)
			return err
		})()
	})))
	t.Run("TimeoutError", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at loading tag by id '%d' from db, case after QueryRow.Scan: %s", 1, "context deadline exceeded")
			_, err := tx.ExecContext(ctx, "SELECT pg_sleep(10)")
			_, err = queries.GetTag(tx, ctx, 1)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
	t.Run("ContextCancelled", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at loading tag by id '%d' from db, case after QueryRow.Scan: %s", 1, "context canceled")
			cancel()
			_, err := queries.GetTag(tx, ctx, 1)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
}

func TestDBTagCreate(t *testing.T) {
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			tagId, err := queries.CreateTag(tx, ctx, TEST_TAG_NAME_1, TEST_TAG_STATE_1)

			assert.Nil(t, err)
			assert.Equal(t, tagId, 1)
			return err
		})()
	})))
	t.Run("DuplicateCase", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			tagId, err := queries.CreateTag(tx, ctx, TEST_TAG_NAME_1, TEST_TAG_STATE_1)

			assert.Nil(t, err)
			assert.NotEqual(t, tagId, -1)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			_, err := queries.CreateTag(tx, ctx, TEST_TAG_NAME_1, TEST_TAG_STATE_1)

			assert.Equal(t, db.ErrorTagDuplicateKey, err)
			return err
		})()
	})))
	t.Run("TimeoutError", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at inserting tag (Name: '%s', State: '%s') into db, case after QueryRow.Scan: %s", TEST_TAG_NAME_1, TEST_TAG_STATE_1, "context deadline exceeded")
			_, err := tx.ExecContext(ctx, "SELECT pg_sleep(10)")
			_, err = queries.CreateTag(tx, ctx, TEST_TAG_NAME_1, TEST_TAG_STATE_1)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
	t.Run("ContextCancelled", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at inserting tag (Name: '%s', State: '%s') into db, case after QueryRow.Scan: %s", TEST_TAG_NAME_1, TEST_TAG_STATE_1, "context canceled")
			cancel()
			_, err := queries.CreateTag(tx, ctx, TEST_TAG_NAME_1, TEST_TAG_STATE_1)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
}

func TestDBTagGetAll(t *testing.T) {
	t.Run("ExpectedEmpty", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			tags, err := queries.GetTags(tx, ctx, 50, 0)

			assert.Nil(t, err)
			assert.Equal(t, 0, len(tags))
			return err
		})()
	})))
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		var expectedTags []entities.Tag
		for i := 1; i <= 10; i++ {
			expectedTags = append(expectedTags, utils.entityGenerators.GenerateTag(i))
		}
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := CreateTagsInDB(t, tx, ctx, 10, TEST_TAG_NAME_TEMPLATE, entities.TAG_STATE_NEW)
			return err
		})()
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			actualTags, err := queries.GetTags(tx, ctx, 50, 0)

			assert.Nil(t, err)
			utils.asserts.AssertEqualTagArrays(t, expectedTags, actualTags)
			return err
		})()
	})))
	t.Run("LimitParameterCase", RunWithRecreateDB((func(t *testing.T) {
		var expectedTags []entities.Tag
		for i := 1; i <= 5; i++ {
			expectedTags = append(expectedTags, utils.entityGenerators.GenerateTag(i))
		}
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := CreateTagsInDB(t, tx, ctx, 10, TEST_TAG_NAME_TEMPLATE, entities.TAG_STATE_NEW)
			return err
		})()
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			actualTags, err := queries.GetTags(tx, ctx, 5, 0)

			assert.Nil(t, err)
			utils.asserts.AssertEqualTagArrays(t, expectedTags, actualTags)
			return err
		})()
	})))
	t.Run("OffsetParameterCase", RunWithRecreateDB((func(t *testing.T) {
		var expectedTags []entities.Tag
		for i := 6; i <= 10; i++ {
			expectedTags = append(expectedTags, utils.entityGenerators.GenerateTag(i))
		}
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := CreateTagsInDB(t, tx, ctx, 10, TEST_TAG_NAME_TEMPLATE, entities.TAG_STATE_NEW)
			return err
		})()
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			actualTags, err := queries.GetTags(tx, ctx, 50, 5)

			assert.Nil(t, err)
			utils.asserts.AssertEqualTagArrays(t, expectedTags, actualTags)
			return err
		})()
	})))
	t.Run("TimeoutError", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at loading tags from db, case after Query: %s", "context deadline exceeded")
			_, err := tx.ExecContext(ctx, "SELECT pg_sleep(10)")
			_, err = queries.GetTags(tx, ctx, 50, 0)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
	t.Run("ContextCancelled", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at loading tags from db, case after Query: %s", "context canceled")
			cancel()
			_, err := queries.GetTags(tx, ctx, 50, 0)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
}

func TestDBTagUpdate(t *testing.T) {
	t.Run("NotFoundCase", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.UpdateTag(tx, ctx, 1, TEST_TAG_NAME_1, TEST_TAG_STATE_1)

			assert.Equal(t, sql.ErrNoRows, err)
			return err
		})()
	})))
	t.Run("DeletedCase", RunWithRecreateDB((func(t *testing.T) {
		expectedTagId := 1
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			tagId, err := queries.CreateTag(tx, ctx, TEST_TAG_NAME_1, TEST_TAG_STATE_1)

			assert.Nil(t, err)
			assert.Equal(t, expectedTagId, tagId)
			return err
		})()
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.DeleteTag(tx, ctx, expectedTagId)

			assert.Nil(t, err)
			return err
		})()
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.UpdateTag(tx, ctx, expectedTagId, TEST_TAG_NAME_2, TEST_TAG_STATE_2)

			assert.Equal(t, sql.ErrNoRows, err)
			return err
		})()
	})))
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		expected := utils.entityGenerators.GenerateTag(1)
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			tagId, err := queries.CreateTag(tx, ctx, TEST_TAG_NAME_2, TEST_TAG_STATE_2)

			assert.Nil(t, err)
			assert.Equal(t, expected.Id, tagId)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.UpdateTag(tx, ctx, expected.Id, expected.Name, expected.State)

			assert.Nil(t, err)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			actual, err := queries.GetTag(tx, ctx, expected.Id)

			utils.asserts.AssertEqualTags(t, expected, actual)
			return err
		})()
	})))
	t.Run("DuplicateCase", RunWithRecreateDB((func(t *testing.T) {
		expectedTagId1 := 1
		expectedTagId2 := 2
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			tagId, err := queries.CreateTag(tx, ctx, TEST_TAG_NAME_1, TEST_TAG_STATE_1)

			assert.Nil(t, err)
			assert.Equal(t, expectedTagId1, tagId)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			tagId, err := queries.CreateTag(tx, ctx, TEST_TAG_NAME_2, TEST_TAG_STATE_2)

			assert.Nil(t, err)
			assert.Equal(t, expectedTagId2, tagId)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.UpdateTag(tx, ctx, expectedTagId2, TEST_TAG_NAME_1, TEST_TAG_STATE_1)

			assert.Equal(t, db.ErrorTagDuplicateKey, err)
			return err
		})()
	})))
	t.Run("TimeoutError", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at updating tag, case after preparing statement: %s", "context deadline exceeded")
			_, err := tx.ExecContext(ctx, "SELECT pg_sleep(10)")
			err = queries.UpdateTag(tx, ctx, 1, TEST_TAG_NAME_1, TEST_TAG_STATE_1)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
	t.Run("ContextCancelled", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at updating tag, case after preparing statement: %s", "context canceled")
			cancel()
			err := queries.UpdateTag(tx, ctx, 1, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
}

func TestDBTagDelete(t *testing.T) {
	t.Run("NotFoundCase", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.DeleteTag(tx, ctx, 1)

			assert.Equal(t, sql.ErrNoRows, err)
			return err
		})()
	})))
	t.Run("AlreadyDeletedCase", RunWithRecreateDB((func(t *testing.T) {
		expectedTagId := 1
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			tagId, err := queries.CreateTag(tx, ctx, TEST_TAG_NAME_1, TEST_TAG_STATE_1)

			assert.Nil(t, err)
			assert.Equal(t, expectedTagId, tagId)
			return err
		})()
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.DeleteTag(tx, ctx, expectedTagId)

			assert.Nil(t, err)
			return err
		})()
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.DeleteTag(tx, ctx, expectedTagId)

			assert.Equal(t, sql.ErrNoRows, err)
			return err
		})()
	})))
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		var expectedTags []entities.Tag
		expectedTags = append(expectedTags, utils.entityGenerators.GenerateTag(1))
		expectedTags = append(expectedTags, utils.entityGenerators.GenerateTag(3))

		tagIdToDelete := 2
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {

			err := CreateTagsInDB(t, tx, ctx, 3, TEST_TAG_NAME_TEMPLATE, entities.TAG_STATE_NEW)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.DeleteTag(tx, ctx, tagIdToDelete)

			assert.Nil(t, err)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			tags, err := queries.GetTags(tx, ctx, 50, 0)

			assert.Nil(t, err)
			utils.asserts.AssertEqualTagArrays(t, expectedTags, tags)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			_, err := queries.GetTag(tx, ctx, tagIdToDelete)

			assert.Equal(t, sql.ErrNoRows, err)
			return err
		})()
	})))
	t.Run("TimeoutError", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at deleting tag, case after preparing statement: %s", "context deadline exceeded")
			_, err := tx.ExecContext(ctx, "SELECT pg_sleep(10)")
			err = queries.DeleteTag(tx, ctx, 1)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
	t.Run("ContextCancelled", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at deleting tag, case after preparing statement: %s", "context canceled")
			cancel()
			err := queries.DeleteTag(tx, ctx, 1)
			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
}
