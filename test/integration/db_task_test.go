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

func TestDBTaskGet(t *testing.T) {
	t.Run("NotFoundCase", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			_, err := queries.GetTask(tx, ctx, 1)

			assert.Equal(t, sql.ErrNoRows, err)
			return err
		})()
	})))
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		expected := utils.entityGenerators.GenerateTask(1)
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			taskId, err := queries.CreateTask(tx, ctx, expected.Name, expected.State)
			assert.Nil(t, err)
			assert.Equal(t, taskId, expected.Id)
			return err
		})()
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			actual, err := queries.GetTask(tx, ctx, expected.Id)
			utils.asserts.AssertEqualTasks(t, expected, actual)
			return err
		})()
	})))
	t.Run("TimeoutError", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at loading task by id '%d' from db, case after QueryRow.Scan: %s", 1, "context deadline exceeded")
			_, err := tx.ExecContext(ctx, "SELECT pg_sleep(10)")
			_, err = queries.GetTask(tx, ctx, 1)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
	t.Run("ContextCancelled", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at loading task by id '%d' from db, case after QueryRow.Scan: %s", 1, "context canceled")
			cancel()
			_, err := queries.GetTask(tx, ctx, 1)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
}

func TestDBTaskCreate(t *testing.T) {
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			taskId, err := queries.CreateTask(tx, ctx, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

			assert.Nil(t, err)
			assert.Equal(t, taskId, 1)
			return err
		})()
	})))
	t.Run("DuplicateCase", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			taskId, err := queries.CreateTask(tx, ctx, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

			assert.Nil(t, err)
			assert.NotEqual(t, taskId, -1)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			_, err := queries.CreateTask(tx, ctx, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

			assert.Equal(t, db.ErrorTaskDuplicateKey, err)
			return err
		})()
	})))
	t.Run("TimeoutError", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at inserting task (Name: '%s', State: '%s') into db, case after QueryRow.Scan: %s", TEST_TASK_NAME_1, TEST_TASK_STATE_1, "context deadline exceeded")
			_, err := tx.ExecContext(ctx, "SELECT pg_sleep(10)")
			_, err = queries.CreateTask(tx, ctx, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
	t.Run("ContextCancelled", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at inserting task (Name: '%s', State: '%s') into db, case after QueryRow.Scan: %s", TEST_TASK_NAME_1, TEST_TASK_STATE_1, "context canceled")
			cancel()
			_, err := queries.CreateTask(tx, ctx, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
}

func TestDBTaskGetAll(t *testing.T) {
	t.Run("ExpectedEmpty", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			tasks, err := queries.GetTasks(tx, ctx, 50, 0)

			assert.Nil(t, err)
			assert.Equal(t, 0, len(tasks))
			return err
		})()
	})))
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		var expectedTasks []entities.Task
		for i := 1; i <= 10; i++ {
			expectedTasks = append(expectedTasks, utils.entityGenerators.GenerateTask(i))
		}
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := CreateTasksInDB(t, tx, ctx, 10, TEST_TASK_NAME_TEMPLATE, entities.TASK_STATE_NEW)
			return err
		})()
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			actualTasks, err := queries.GetTasks(tx, ctx, 50, 0)

			assert.Nil(t, err)
			utils.asserts.AssertEqualTaskArrays(t, expectedTasks, actualTasks)
			return err
		})()
	})))
	t.Run("LimitParameterCase", RunWithRecreateDB((func(t *testing.T) {
		var expectedTasks []entities.Task
		for i := 1; i <= 5; i++ {
			expectedTasks = append(expectedTasks, utils.entityGenerators.GenerateTask(i))
		}
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := CreateTasksInDB(t, tx, ctx, 10, TEST_TASK_NAME_TEMPLATE, entities.TASK_STATE_NEW)
			return err
		})()
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			actualTasks, err := queries.GetTasks(tx, ctx, 5, 0)

			assert.Nil(t, err)
			utils.asserts.AssertEqualTaskArrays(t, expectedTasks, actualTasks)
			return err
		})()
	})))
	t.Run("OffsetParameterCase", RunWithRecreateDB((func(t *testing.T) {
		var expectedTasks []entities.Task
		for i := 6; i <= 10; i++ {
			expectedTasks = append(expectedTasks, utils.entityGenerators.GenerateTask(i))
		}
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := CreateTasksInDB(t, tx, ctx, 10, TEST_TASK_NAME_TEMPLATE, entities.TASK_STATE_NEW)
			return err
		})()
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			actualTasks, err := queries.GetTasks(tx, ctx, 50, 5)

			assert.Nil(t, err)
			utils.asserts.AssertEqualTaskArrays(t, expectedTasks, actualTasks)
			return err
		})()
	})))
	t.Run("TimeoutError", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at loading tasks from db, case after Query: %s", "context deadline exceeded")
			_, err := tx.ExecContext(ctx, "SELECT pg_sleep(10)")
			_, err = queries.GetTasks(tx, ctx, 50, 0)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
	t.Run("ContextCancelled", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at loading tasks from db, case after Query: %s", "context canceled")
			cancel()
			_, err := queries.GetTasks(tx, ctx, 50, 0)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
}

func TestDBTaskUpdate(t *testing.T) {
	t.Run("NotFoundCase", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.UpdateTask(tx, ctx, 1, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

			assert.Equal(t, sql.ErrNoRows, err)
			return err
		})()
	})))
	t.Run("DeletedCase", RunWithRecreateDB((func(t *testing.T) {
		expectedTaskId := 1
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			taskId, err := queries.CreateTask(tx, ctx, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

			assert.Nil(t, err)
			assert.Equal(t, expectedTaskId, taskId)
			return err
		})()
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.DeleteTask(tx, ctx, expectedTaskId)

			assert.Nil(t, err)
			return err
		})()
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.UpdateTask(tx, ctx, expectedTaskId, TEST_TASK_NAME_2, TEST_TASK_STATE_2)

			assert.Equal(t, sql.ErrNoRows, err)
			return err
		})()
	})))
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		expected := utils.entityGenerators.GenerateTask(1)
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			taskId, err := queries.CreateTask(tx, ctx, TEST_TASK_NAME_2, TEST_TASK_STATE_2)

			assert.Nil(t, err)
			assert.Equal(t, expected.Id, taskId)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.UpdateTask(tx, ctx, expected.Id, expected.Name, expected.State)

			assert.Nil(t, err)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			actual, err := queries.GetTask(tx, ctx, expected.Id)

			utils.asserts.AssertEqualTasks(t, expected, actual)
			return err
		})()
	})))
	t.Run("DuplicateCase", RunWithRecreateDB((func(t *testing.T) {
		expectedTaskId1 := 1
		expectedTaskId2 := 2
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			taskId, err := queries.CreateTask(tx, ctx, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

			assert.Nil(t, err)
			assert.Equal(t, expectedTaskId1, taskId)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			taskId, err := queries.CreateTask(tx, ctx, TEST_TASK_NAME_2, TEST_TASK_STATE_2)

			assert.Nil(t, err)
			assert.Equal(t, expectedTaskId2, taskId)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.UpdateTask(tx, ctx, expectedTaskId2, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

			assert.Equal(t, db.ErrorTaskDuplicateKey, err)
			return err
		})()
	})))
	t.Run("TimeoutError", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at updating task, case after preparing statement: %s", "context deadline exceeded")
			_, err := tx.ExecContext(ctx, "SELECT pg_sleep(10)")
			err = queries.UpdateTask(tx, ctx, 1, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
	t.Run("ContextCancelled", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at updating task, case after preparing statement: %s", "context canceled")
			cancel()
			err := queries.UpdateTask(tx, ctx, 1, TEST_TASK_NAME_1, TEST_TASK_STATE_1)
			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
}

func TestDBTaskDelete(t *testing.T) {
	t.Run("NotFoundCase", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.DeleteTask(tx, ctx, 1)

			assert.Equal(t, sql.ErrNoRows, err)
			return err
		})()
	})))
	t.Run("AlreadyDeletedCase", RunWithRecreateDB((func(t *testing.T) {
		expectedTaskId := 1
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			taskId, err := queries.CreateTask(tx, ctx, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

			assert.Nil(t, err)
			assert.Equal(t, expectedTaskId, taskId)
			return err
		})()
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.DeleteTask(tx, ctx, expectedTaskId)

			assert.Nil(t, err)
			return err
		})()
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.DeleteTask(tx, ctx, expectedTaskId)

			assert.Equal(t, sql.ErrNoRows, err)
			return err
		})()
	})))
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		var expectedTasks []entities.Task
		expectedTasks = append(expectedTasks, utils.entityGenerators.GenerateTask(1))
		expectedTasks = append(expectedTasks, utils.entityGenerators.GenerateTask(3))

		taskIdToDelete := 2
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {

			err := CreateTasksInDB(t, tx, ctx, 3, TEST_TASK_NAME_TEMPLATE, entities.TASK_STATE_NEW)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.DeleteTask(tx, ctx, taskIdToDelete)

			assert.Nil(t, err)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			tasks, err := queries.GetTasks(tx, ctx, 50, 0)

			assert.Nil(t, err)
			utils.asserts.AssertEqualTaskArrays(t, expectedTasks, tasks)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			_, err := queries.GetTask(tx, ctx, taskIdToDelete)

			assert.Equal(t, sql.ErrNoRows, err)
			return err
		})()
	})))
	t.Run("TimeoutError", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at deleting task, case after preparing statement: %s", "context deadline exceeded")
			_, err := tx.ExecContext(ctx, "SELECT pg_sleep(10)")
			err = queries.DeleteTask(tx, ctx, 1)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
	t.Run("ContextCancelled", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at deleting task, case after preparing statement: %s", "context canceled")
			cancel()
			err := queries.DeleteTask(tx, ctx, 1)
			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
}
