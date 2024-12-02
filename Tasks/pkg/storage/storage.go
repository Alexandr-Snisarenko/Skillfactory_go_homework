package storage

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Хранилище данных.
type Storage struct {
	db *pgxpool.Pool
}

// Конструктор, принимает строку подключения к БД.
func NewStorage(dbURL string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		return nil, err
	}
	s := Storage{
		db: db,
	}
	return &s, nil
}

// Задача.
type Task struct {
	ID         int
	Opened     int64
	Closed     int64
	AuthorID   int
	AssignedID int
	Title      string
	Content    string
}

// NewTask создаёт новую задачу и возвращает её id.
func (s *Storage) NewTask(title, content string) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO tasks (title, content)
		VALUES ($1, $2) RETURNING id;
		`,
		title,
		content,
	).Scan(&id)
	return id, err
}

// UpdateTask обновляет задачу по id.
func (s *Storage) UpdateTask(t *Task) error {
	_, err := s.db.Exec(context.Background(), `
		UPDATE tasks 
		SET opened = $1,
			closed = $2,
			author_id = $3, 
			assigned_id = $4, 
			title = $5, 
			content = $6
		WHERE id = $7;
		`,
		t.Opened,
		t.Closed,
		t.AuthorID,
		t.AssignedID,
		t.Title,
		t.Content,
		t.ID,
	)

	return err
}

// DeleteTask удаляет задачу по id.
func (s *Storage) DeleteTask(t *Task) error {
	_, err := s.db.Exec(context.Background(), `
		DELETE FROM tasks 
		WHERE id = $1;
		`,
		t.ID,
	)

	return err
}

// GetTask возвращает задачу из БД по ID.
func (s *Storage) GetTask(taskID int) (*Task, error) {
	t := new(Task)

	err := s.db.QueryRow(context.Background(), `
	  SELECT 
			id,
			opened,
			closed,
			author_id,
			assigned_id,
			title,
			content
		FROM tasks
	   WHERE id = $1
		ORDER BY id;
	`,
		taskID,
	).Scan(&t.ID, &t.Opened, &t.Closed, &t.AuthorID, &t.AssignedID, &t.Title, &t.Content)

	if err != nil {
		return nil, err
	}

	return t, nil
}

// Tasks возвращает список задач по фильтру (assignedID, authorID) из БД.
func (s *Storage) Tasks(assignedID, authorID int) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
	  SELECT 
			id,
			opened,
			closed,
			author_id,
			assigned_id,
			title,
			content
		FROM tasks
	   WHERE 
			($1 = 0 OR assigned_id = $1) AND
			($2 = 0 OR author_id = $2)	   
		ORDER BY id;
	`,
		assignedID,
		authorID,
	)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		tasks = append(tasks, t)

	}

	// ВАЖНО не забыть проверить rows.Err()
	return tasks, rows.Err()
}
