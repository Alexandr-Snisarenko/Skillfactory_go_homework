package mongo

import (
	"GoNews/pkg/storage"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	databaseName   = "GoNews" // имя учебной БД
	collectionName = "posts"  // имя коллекции в учебной БД
)

// Хранилище данных.
type Store struct {
	collection *mongo.Collection
}

// Конструктор объекта хранилища.
func NewStorage(uri string) (*Store, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Проверяем соединение
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return &Store{
		collection: client.Database(databaseName).Collection(collectionName),
	}, nil
}

// Добавление нового поста
func (ms *Store) AddPost(post storage.Post) error {
	// при добавлении используем INT ID в соответствии со структурой post из ТЗ
	// для этого генерим псевдоуникальный ID на сонове UNIX микросекунд
	post.ID = time.Now().UnixMicro()
	_, err := ms.collection.InsertOne(context.Background(), post)
	return err
}

// Обновление поста
func (ms *Store) UpdatePost(post storage.Post) error {
	filter := bson.M{"id": post.ID}
	updates := bson.M{"$set": bson.M{"title": post.Title, "content": post.Content}}
	_, err := ms.collection.UpdateOne(context.Background(), filter, updates)
	return err
}

// Список всех постов
func (ms *Store) Posts() ([]storage.Post, error) {
	var posts []storage.Post
	cursor, err := ms.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var post storage.Post
		err := cursor.Decode(&post)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// Удаление поста
func (ms *Store) DeletePost(post storage.Post) error {
	filter := bson.M{"id": post.ID}
	_, err := ms.collection.DeleteOne(context.Background(), filter)
	return err
}

// Поиск  поста по id
// Предполагаем, что id - уникальный
func (ms *Store) GetPost(id int64) (*storage.Post, error) {
	var posts []storage.Post
	cursor, err := ms.collection.Find(context.Background(), bson.M{"id": id})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var post storage.Post
		err := cursor.Decode(&post)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return &posts[0], nil
}
