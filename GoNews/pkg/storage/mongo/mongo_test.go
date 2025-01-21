package mongo

import (
	"GoNews/pkg/storage"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMongoCRUD(t *testing.T) {
	var (
		dbURL string = "mongodb://localhost:27017/"
		db    *Store
		posts []storage.Post
		post  storage.Post = storage.Post{
			ID:         0,
			Title:      "Test Title",
			Content:    "Test Content",
			AuthorID:   0,
			AuthorName: "test author",
			CreatedAt:  0,
		}
		postCnt int
	)

	t.Run("RunTest", func(t *testing.T) {
		var err error
		t.Log("Connect DB")
		db, err = NewStorage(dbURL)
		require.NoError(t, err)

		// Get Posts
		t.Log("Get Posts")
		posts, err = db.Posts()
		require.NoError(t, err)
		// сохраняем текущее кол-во записей
		postCnt = len(posts)

		// Insert Post
		t.Log("Insert Post")
		err = db.AddPost(post)
		require.NoError(t, err)

		// проверяем, что запись добавилась
		posts, err = db.Posts()
		require.NoError(t, err)
		require.NotEmpty(t, posts)
		//проверяем, что записей стало на 1 больше
		require.True(t, (len(posts) == postCnt+1))
		postCnt = len(posts)
		post = posts[postCnt-1]

		// Update Post
		t.Log("Update Post")
		post.Content = "check update"
		err = db.UpdatePost(post)
		require.NoError(t, err)

		// check Update
		ppost, err := db.GetPost(post.ID)
		require.NoError(t, err)
		require.NotEmpty(t, ppost)
		//проверяем, что обновилась запись
		post = *ppost
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
