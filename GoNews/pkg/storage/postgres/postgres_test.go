package postgres

import (
	"GoNews/pkg/storage"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPostgresCRUD(t *testing.T) {
	var (
		dbURL   string = "postgres://postgres:postgres@localhost/postgres"
		db      *Store
		posts   []storage.Post
		post    storage.Post
		postCnt int
	)

	// для работы тестов в БД должна быть хоть одна запись
	// ввиду того, что в текущем варианте функциональности добавления авторов нет
	// то лгика теста следующая:
	// получаем все записи post из БД
	// берём крайний объект post из полученного списка, вносим в него изменения,
	// создаём новую запись на его основе, проверяем, что запись добавилась,
	// вносим в неё изменения, проверяем, что изменения применились,
	// удаляем добавленную запись, проверяем, что запись удалтлась.
	t.Run("RunTest", func(t *testing.T) {
		var err error
		t.Log("Connect DB")
		db, err = NewStorage(dbURL)
		require.NoError(t, err)

		// Get Posts
		t.Log("Get Posts")
		posts, err = db.Posts()

		require.NoError(t, err)
		require.NotEmpty(t, posts)

		// сохраняем текущее колво записей
		postCnt = len(posts)
		//берем крайний post из записей
		post = posts[postCnt-1]

		// Insert Post
		t.Log("Insert Post")
		post.Title = "Test"
		post.Content = "Test"
		err = db.AddPost(post)
		require.NoError(t, err)

		// проверяем, что запись добавилась
		posts, err = db.Posts()

		require.NoError(t, err)
		require.NotEmpty(t, posts)
		//проверяем, что добавилсась запись
		require.True(t, (len(posts) == postCnt+1))
		postCnt = len(posts)
		post = posts[postCnt-1]

		// Update Post
		t.Log("Update Post")
		post.Content = "check update"
		err = db.UpdatePost(post)
		require.NoError(t, err)

		// check Update
		posts, err = db.Posts()

		require.NoError(t, err)
		require.NotEmpty(t, posts)
		//проверяем, что обновилась запись
		post = posts[postCnt-1]
		require.EqualValues(t, "check update", post.Content)

		// Delete Post
		t.Log("Delete Post")
		err = db.DeletePost(post)
		require.NoError(t, err)

		// check Delete
		posts, err = db.Posts()

		require.NoError(t, err)
		require.NotEmpty(t, posts)
		//проверяем, что удалилась запись
		require.True(t, (len(posts) == postCnt-1))
	})

}
