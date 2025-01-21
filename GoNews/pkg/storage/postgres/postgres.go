package postgres

import (
	"GoNews/pkg/storage"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Хранилище данных.
type Store struct {
	db *pgxpool.Pool
}

// Конструктор, принимает строку подключения к БД.
func NewStorage(dbURL string) (*Store, error) {
	db, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Posts() ([]storage.Post, error) {
	query := `
	SELECT p.id, p.title, p.content, a.id, a.name, p.created_at 
	FROM authors a, posts p
	WHERE p.author_id = a.id
	`
	rows, err := s.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []storage.Post
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var p storage.Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.AuthorID,
			&p.AuthorName,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		// ошибка при получении результатов
		return nil, err
	}

	return posts, rows.Err()
}

func (s *Store) AddPost(p storage.Post) error {
	query := ` 
	INSERT INTO posts (title, content, author_id, created_at)
	VALUES ($1, $2, $3, $4)
	`
	_, err := s.db.Exec(context.Background(), query,
		p.Title,
		p.Content,
		p.AuthorID,
		p.CreatedAt,
	)
	return err
}

func (s *Store) UpdatePost(p storage.Post) error {
	query := ` 
	UPDATE posts
	SET  title = $1,
	content = $2, 
	author_id = $3, 
	created_at = $4
	WHERE id = $5
	`
	_, err := s.db.Exec(context.Background(), query,
		p.Title,
		p.Content,
		p.AuthorID,
		p.CreatedAt,
		p.ID,
	)
	return err
}

func (s *Store) DeletePost(p storage.Post) error {
	query := "DELETE FROM posts WHERE id = $1"
	_, err := s.db.Exec(context.Background(), query, p.ID)
	return err
}
