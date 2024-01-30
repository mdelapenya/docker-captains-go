package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type todoStore interface {
	Create(ctx context.Context, t *Todo) error
	DeleteAll(ctx context.Context) error
	DeleteByID(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (Todo, error)
	List(ctx context.Context) ([]Todo, error)
	Update(ctx context.Context, t Todo) error
}

type Repository struct {
	conn *pgx.Conn
}

// NewTodosRepository creates a new repository. It will receive a context and the PostgreSQL connection string.
func NewTodosRepository(ctx context.Context, connStr string) (*Repository, error) {
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
	}

	return &Repository{
		conn: conn,
	}, nil
}

// Create creates a new taTodok in the database.
// It uses value semantics at the method receiver to avoid mutating the original repository.
// It uses pointer semantics at the todo parameter to avoid copying the struct, modifying it and returning it.
func (r Repository) Create(ctx context.Context, t *Todo) error {
	query := "INSERT INTO todos (id, title, completed, order) VALUES ($1, $2, $3, $4) returning id"

	return r.conn.QueryRow(ctx, query, t.ID, t.Title, t.Completed, t.Order).Scan(&t.ID)
}

// DeleteByID deletes a todo from the database by its ID.
func (r Repository) DeleteByID(ctx context.Context, id string) error {
	query := "DELETE FROM todos WHERE id = $1"

	return r.conn.QueryRow(ctx, query, id).Scan()
}

// DeleteAll deletes all todo from the database.
func (r Repository) DeleteAll(ctx context.Context) error {
	query := "DELETE FROM todos"

	return r.conn.QueryRow(ctx, query).Scan()
}

// FindByID retrieves a todo from the database by its ID.
func (r Repository) FindByID(ctx context.Context, id string) (Todo, error) {
	query := "SELECT id, title, completed, order FROM todos WHERE id = $1"

	var t Todo
	err := r.conn.QueryRow(ctx, query, id).Scan(&t.ID, &t.Title, &t.Completed, &t.Order)
	if err != nil {
		return Todo{}, err
	}

	return t, nil
}

// List retrieves all todo from the database (no filters nor pagination).
func (r Repository) List(ctx context.Context) ([]Todo, error) {
	query := "SELECT id, title, completed, order FROM todos"

	var ts []Todo
	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return []Todo{}, err
	}

	err = rows.Scan(ts)
	if err != nil {
		return []Todo{}, err
	}

	return ts, nil
}

// Update updates a todo in the database. The todo is identified by its ID, and
// the new values are taken from the todo parameter.
func (r Repository) Update(ctx context.Context, t Todo) error {
	existingTodo, err := r.FindByID(ctx, t.ID)
	if err != nil {
		return fmt.Errorf("todo not found: %w", err)
	}

	existingTodo.Title = t.Title
	existingTodo.Completed = t.Completed
	existingTodo.Order = t.Order

	query := "UPDATE todos SET title = $2, completed = $3, order = $4 WHERE id = $1"

	return r.conn.QueryRow(ctx, query, t.ID, existingTodo.Title, existingTodo.Completed, existingTodo.Order).Scan()
}
