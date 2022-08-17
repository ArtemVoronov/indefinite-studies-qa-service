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

func TestDBNoteGet(t *testing.T) {
	t.Run("NotFoundCase", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			_, err := queries.GetNote(tx, ctx, 1)

			assert.Equal(t, sql.ErrNoRows, err)
			return err
		})()
	})))
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		result, _ := db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
			userId, err := CreateUserInDB(t, tx, ctx, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
			return userId, err
		})()
		userId, ok := result.(int)
		assert.True(t, ok)

		result, _ = db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
			tagId, err := CreateTagInDB(t, tx, ctx, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
			return tagId, err
		})()
		tagId, ok := result.(int)
		assert.True(t, ok)

		expected := utils.entityGenerators.GenerateNote(1, userId, tagId)

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			noteId, err := queries.CreateNote(tx, ctx, expected.Text, expected.Topic, expected.TagId, expected.UserId, expected.State)

			assert.Nil(t, err)
			assert.Equal(t, expected.Id, noteId)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			actual, err := queries.GetNote(tx, ctx, expected.Id)

			utils.asserts.AssertEqualNotes(t, expected, actual)
			return err
		})()
	})))
	t.Run("TimeoutError", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at loading note by id '%d' from db, case after QueryRow.Scan: %s", 1, "context deadline exceeded")
			_, err := tx.ExecContext(ctx, "SELECT pg_sleep(10)")
			_, err = queries.GetNote(tx, ctx, 1)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
	t.Run("ContextCancelled", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at loading note by id '%d' from db, case after QueryRow.Scan: %s", 1, "context canceled")
			cancel()
			_, err := queries.GetNote(tx, ctx, 1)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
}

func TestDBNoteCreate(t *testing.T) {
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		result, _ := db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
			userId, err := CreateUserInDB(t, tx, ctx, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
			return userId, err
		})()
		userId, ok := result.(int)
		assert.True(t, ok)

		result, _ = db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
			tagId, err := CreateTagInDB(t, tx, ctx, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
			return tagId, err
		})()
		tagId, ok := result.(int)
		assert.True(t, ok)

		expected := utils.entityGenerators.GenerateNote(1, userId, tagId)

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			noteId, err := queries.CreateNote(tx, ctx, expected.Text, expected.Topic, expected.TagId, expected.UserId, expected.State)

			assert.Nil(t, err)
			assert.Equal(t, expected.Id, noteId)
			return err
		})()
	})))
	t.Run("TimeoutError", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at inserting note (Topic: '%s', UserId: '%d') into db, case after QueryRow.Scan: %s", TEST_NOTE_TOPIC_1, 1, "context deadline exceeded")
			_, err := tx.ExecContext(ctx, "SELECT pg_sleep(10)")
			_, err = queries.CreateNote(tx, ctx, TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, 1, 1, TEST_NOTE_STATE_1)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
	t.Run("ContextCancelled", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at inserting note (Topic: '%s', UserId: '%d') into db, case after QueryRow.Scan: %s", TEST_NOTE_TOPIC_1, 1, "context canceled")
			cancel()
			_, err := queries.CreateNote(tx, ctx, TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, 1, 1, TEST_NOTE_STATE_1)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
}

func TestDBNoteGetAll(t *testing.T) {
	t.Run("ExpectedEmpty", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			notes, err := queries.GetNotes(tx, ctx, 50, 0)

			assert.Nil(t, err)
			assert.Equal(t, 0, len(notes))
			return err
		})()
	})))
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		result, _ := db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
			userId, err := CreateUserInDB(t, tx, ctx, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
			return userId, err
		})()
		userId, ok := result.(int)
		assert.True(t, ok)

		result, _ = db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
			tagId, err := CreateTagInDB(t, tx, ctx, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
			return tagId, err
		})()
		tagId, ok := result.(int)
		assert.True(t, ok)

		var expectedNotes []entities.Note
		for i := 1; i <= 10; i++ {
			expectedNotes = append(expectedNotes, utils.entityGenerators.GenerateNote(i, userId, tagId))
		}

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := CreateNotesInDB(t, tx, ctx, 10, TEST_NOTE_TEXT_TEMPLATE, TEST_NOTE_TOPIC_TEMPLATE, tagId, userId, TEST_NOTE_STATE_1)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			actualNotes, err := queries.GetNotes(tx, ctx, 50, 0)

			assert.Nil(t, err)
			utils.asserts.AssertEqualNoteArrays(t, expectedNotes, actualNotes)
			return err
		})()
	})))
	t.Run("LimitParameterCase", RunWithRecreateDB((func(t *testing.T) {
		result, _ := db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
			userId, err := CreateUserInDB(t, tx, ctx, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
			return userId, err
		})()
		userId, ok := result.(int)
		assert.True(t, ok)

		result, _ = db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
			tagId, err := CreateTagInDB(t, tx, ctx, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
			return tagId, err
		})()
		tagId, ok := result.(int)
		assert.True(t, ok)

		var expectedNotes []entities.Note
		for i := 1; i <= 5; i++ {
			expectedNotes = append(expectedNotes, utils.entityGenerators.GenerateNote(i, userId, tagId))
		}

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := CreateNotesInDB(t, tx, ctx, 10, TEST_NOTE_TEXT_TEMPLATE, TEST_NOTE_TOPIC_TEMPLATE, tagId, userId, TEST_NOTE_STATE_1)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			actualNotes, err := queries.GetNotes(tx, ctx, 5, 0)

			assert.Nil(t, err)
			utils.asserts.AssertEqualNoteArrays(t, expectedNotes, actualNotes)
			return err
		})()
	})))
	t.Run("OffsetParameterCase", RunWithRecreateDB((func(t *testing.T) {
		result, _ := db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
			userId, err := CreateUserInDB(t, tx, ctx, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
			return userId, err
		})()
		userId, ok := result.(int)
		assert.True(t, ok)

		result, _ = db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
			tagId, err := CreateTagInDB(t, tx, ctx, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
			return tagId, err
		})()
		tagId, ok := result.(int)
		assert.True(t, ok)

		var expectedNotes []entities.Note
		for i := 6; i <= 10; i++ {
			expectedNotes = append(expectedNotes, utils.entityGenerators.GenerateNote(i, userId, tagId))
		}

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := CreateNotesInDB(t, tx, ctx, 10, TEST_NOTE_TEXT_TEMPLATE, TEST_NOTE_TOPIC_TEMPLATE, tagId, userId, TEST_NOTE_STATE_1)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			actualNotes, err := queries.GetNotes(tx, ctx, 50, 5)

			assert.Nil(t, err)
			utils.asserts.AssertEqualNoteArrays(t, expectedNotes, actualNotes)
			return err
		})()
	})))
	t.Run("TimeoutError", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at loading note from db, case after Query: %s", "context deadline exceeded")
			_, err := tx.ExecContext(ctx, "SELECT pg_sleep(10)")
			_, err = queries.GetNotes(tx, ctx, 50, 0)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
	t.Run("ContextCancelled", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at loading note from db, case after Query: %s", "context canceled")
			cancel()
			_, err := queries.GetNotes(tx, ctx, 50, 0)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
}

func TestDBNoteUpdate(t *testing.T) {
	t.Run("NotFoundCase", RunWithRecreateDB((func(t *testing.T) {
		result, err := db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
			userId, err := CreateUserInDB(t, tx, ctx, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
			return userId, err
		})()
		userId, ok := result.(int)
		assert.True(t, ok)

		result, err = db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
			tagId, err := CreateTagInDB(t, tx, ctx, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
			return tagId, err
		})()
		tagId, ok := result.(int)
		assert.True(t, ok)

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err = queries.UpdateNote(tx, ctx, 1, TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, tagId, userId, TEST_NOTE_STATE_1)

			assert.Equal(t, sql.ErrNoRows, err)
			return err
		})()
	})))
	t.Run("DeletedCase", RunWithRecreateDB((func(t *testing.T) {
		expectedNoteId := 1
		result, err := db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
			userId, err := CreateUserInDB(t, tx, ctx, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
			return userId, err
		})()
		userId, ok := result.(int)
		assert.True(t, ok)

		result, err = db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
			tagId, err := CreateTagInDB(t, tx, ctx, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
			return tagId, err
		})()
		tagId, ok := result.(int)
		assert.True(t, ok)

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			noteId, err := queries.CreateNote(tx, ctx, TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, tagId, userId, TEST_NOTE_STATE_1)

			assert.Nil(t, err)
			assert.Equal(t, expectedNoteId, noteId)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err = queries.DeleteNote(tx, ctx, expectedNoteId)

			assert.Nil(t, err)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err = queries.UpdateNote(tx, ctx, expectedNoteId, TEST_NOTE_TEXT_2, TEST_NOTE_TOPIC_2, tagId, userId, TEST_NOTE_STATE_2)

			assert.Equal(t, sql.ErrNoRows, err)
			return err
		})()
	})))
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		result, err := db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
			userId, err := CreateUserInDB(t, tx, ctx, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
			return userId, err
		})()
		userId, ok := result.(int)
		assert.True(t, ok)

		result, err = db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
			tagId, err := CreateTagInDB(t, tx, ctx, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
			return tagId, err
		})()
		tagId, ok := result.(int)
		assert.True(t, ok)

		expected := utils.entityGenerators.GenerateNote(1, userId, tagId)

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			noteId, err := queries.CreateNote(tx, ctx, TEST_NOTE_TEXT_2, TEST_NOTE_TOPIC_2, tagId, userId, TEST_NOTE_STATE_2)

			assert.Nil(t, err)
			assert.Equal(t, expected.Id, noteId)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err = queries.UpdateNote(tx, ctx, expected.Id, expected.Text, expected.Topic, expected.TagId, expected.UserId, expected.State)

			assert.Nil(t, err)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			actual, err := queries.GetNote(tx, ctx, expected.Id)

			utils.asserts.AssertEqualNotes(t, expected, actual)
			return err
		})()
	})))
	t.Run("TimeoutError", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at updating note, case after preparing statement: %s", "context deadline exceeded")
			_, err := tx.ExecContext(ctx, "SELECT pg_sleep(10)")
			err = queries.UpdateNote(tx, ctx, 1, TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, 1, 1, TEST_NOTE_STATE_1)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
	t.Run("ContextCancelled", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at updating note, case after preparing statement: %s", "context canceled")
			cancel()
			err := queries.UpdateNote(tx, ctx, 1, TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, 1, 1, TEST_NOTE_STATE_1)
			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
}

func TestDBNoteDelete(t *testing.T) {
	t.Run("NotFoundCase", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.DeleteNote(tx, ctx, 1)

			assert.Equal(t, sql.ErrNoRows, err)
			return err
		})()
	})))
	t.Run("AlreadyDeletedCase", RunWithRecreateDB((func(t *testing.T) {
		result, err := db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
			userId, err := CreateUserInDB(t, tx, ctx, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
			return userId, err
		})()
		userId, ok := result.(int)
		assert.True(t, ok)

		result, err = db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
			tagId, err := CreateTagInDB(t, tx, ctx, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
			return tagId, err
		})()
		tagId, ok := result.(int)
		assert.True(t, ok)

		expectedNoteId := 1

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			noteId, err := queries.CreateNote(tx, ctx, TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, tagId, userId, TEST_NOTE_STATE_1)

			assert.Nil(t, err)
			assert.Equal(t, expectedNoteId, noteId)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err = queries.DeleteNote(tx, ctx, expectedNoteId)

			assert.Nil(t, err)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err = queries.DeleteNote(tx, ctx, expectedNoteId)

			assert.Equal(t, sql.ErrNoRows, err)
			return err
		})()
	})))
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		result, err := db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
			userId, err := CreateUserInDB(t, tx, ctx, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
			return userId, err
		})()
		userId, ok := result.(int)
		assert.True(t, ok)

		result, err = db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
			tagId, err := CreateTagInDB(t, tx, ctx, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
			return tagId, err
		})()
		tagId, ok := result.(int)
		assert.True(t, ok)

		var expectedNotes []entities.Note
		expectedNotes = append(expectedNotes, utils.entityGenerators.GenerateNote(1, userId, tagId))
		expectedNotes = append(expectedNotes, utils.entityGenerators.GenerateNote(3, userId, tagId))

		noteIdToDelete := 2

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := CreateNotesInDB(t, tx, ctx, 3, TEST_NOTE_TEXT_TEMPLATE, TEST_NOTE_TOPIC_TEMPLATE, tagId, userId, TEST_NOTE_STATE_1)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err = queries.DeleteNote(tx, ctx, noteIdToDelete)

			assert.Nil(t, err)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			notes, err := queries.GetNotes(tx, ctx, 50, 0)

			assert.Nil(t, err)
			utils.asserts.AssertEqualNoteArrays(t, expectedNotes, notes)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			_, err = queries.GetNote(tx, ctx, noteIdToDelete)

			assert.Equal(t, sql.ErrNoRows, err)
			return err
		})()
	})))
	t.Run("TimeoutError", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at deleting note, case after preparing statement: %s", "context deadline exceeded")
			_, err := tx.ExecContext(ctx, "SELECT pg_sleep(10)")
			err = queries.DeleteNote(tx, ctx, 1)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
	t.Run("ContextCancelled", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at deleting note, case after preparing statement: %s", "context canceled")
			cancel()
			err := queries.DeleteNote(tx, ctx, 1)
			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
}
