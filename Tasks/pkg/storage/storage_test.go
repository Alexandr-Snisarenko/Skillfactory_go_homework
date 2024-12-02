package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCRUD(t *testing.T) {
	var (
		dbURL  string = "postgres://postgres:postgres@localhost/postgres"
		db     *Storage
		taskID int
		task   *Task
		tasks  []Task
	)

	t.Run("coonect DB", func(t *testing.T) {
		var err error
		db, err = NewStorage(dbURL)

		require.NoError(t, err)
	})

	t.Run("Insert Task", func(t *testing.T) {
		var err error
		taskID, err = db.NewTask("test", "test task")

		require.NoError(t, err)
		require.NotEmpty(t, taskID)
	})

	t.Run("Get Task", func(t *testing.T) {
		var err error
		task, err = db.GetTask(taskID)

		require.NoError(t, err)
		require.NotEmpty(t, task)
	})

	t.Run("Select Tasks", func(t *testing.T) {
		var err error
		tasks, err = db.Tasks(0, 0)

		require.NoError(t, err)
		require.NotEmpty(t, tasks)
	})

	t.Run("Update Task", func(t *testing.T) {
		var err error
		task.Content = "check update"

		err = db.UpdateTask(task)
		require.NoError(t, err)

		// check Update
		task, err = db.GetTask(taskID)
		require.NoError(t, err)
		require.EqualValues(t, "check update", task.Content)

	})

	t.Run("Delete Task", func(t *testing.T) {
		var err error

		err = db.DeleteTask(task)
		require.NoError(t, err)

		// check Delete
		task, err = db.GetTask(taskID)
		require.Error(t, err)
		require.EqualError(t, err, "no rows in result set")
	})

}
