package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
	pool *pgxpool.Pool
}

func dbConfig(connStr string) *pgxpool.Config {
	const defaultMaxConns = int32(4)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	// Your own Database URL
	var DATABASE_URL string = connStr

	dbConfig, err := pgxpool.ParseConfig(DATABASE_URL)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		return true
	}

	dbConfig.BeforeClose = func(c *pgx.Conn) {
	}

	return dbConfig
}

var ErrNotFound = fmt.Errorf("not found")

// NewTodosRepository creates a new repository. It will receive a context and the PostgreSQL connection string.
func NewTodosRepository(ctx context.Context, connStr string) (*Repository, error) {
	connPool, err := pgxpool.NewWithConfig(context.Background(), dbConfig(connStr))
	if err != nil {
		log.Fatal("Error while creating connection to the database!!")
	}

	connection, err := connPool.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error while acquiring connection from the database pool: %w", err)
	}
	defer connection.Release()

	err = connection.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not ping database: %w", err)
	}

	return &Repository{
		pool: connPool,
	}, nil
}

// Create creates a new todo in the database.
// It uses value semantics at the method receiver to avoid mutating the original repository.
// It uses pointer semantics at the todo parameter to avoid copying the struct, modifying it and returning it.
func (r Repository) Create(ctx context.Context, t *Todo) error {
	query := "INSERT INTO todos (id, title, completed, order_number) VALUES ($1, $2, $3, $4) returning id"

	if t.ID == "" {
		t.ID = uuid.NewString()
	}

	err := r.pool.QueryRow(ctx, query, t.ID, t.Title, t.Completed, t.Order).Scan(&t.ID)
	if err != nil {
		return err
	}

	t.Url = t.URL()

	return nil
}

// DeleteByID deletes a todo from the database by its ID.
func (r Repository) DeleteByID(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM todos WHERE id = $1", id)
	return err
}

// DeleteAll deletes all todo from the database.
func (r Repository) DeleteAll(ctx context.Context) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM todos")
	return err
}

// FindByID retrieves a todo from the database by its ID.
func (r Repository) FindByID(ctx context.Context, id string) (Todo, error) {
	query := "SELECT id, title, completed, order_number FROM todos WHERE id = $1"

	var t Todo
	err := r.pool.QueryRow(ctx, query, id).Scan(&t.ID, &t.Title, &t.Completed, &t.Order)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Todo{}, ErrNotFound
		}

		return Todo{}, err
	}

	return t, nil
}

// List retrieves all todos from the database (no filters nor pagination).
func (r Repository) List(ctx context.Context) ([]Todo, error) {
	query := "SELECT id, title, completed, order_number FROM todos"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return []Todo{}, err
	}

	var ts []Todo
	for rows.Next() {
		var t Todo
		err := rows.Scan(&t.ID, &t.Title, &t.Completed, &t.Order)
		if err != nil {
			return ts, err
		}

		t.Url = t.URL()

		ts = append(ts, t)
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

	query := "UPDATE todos SET title = $2, completed = $3, order_number = $4 WHERE id = $1"

	err = r.pool.QueryRow(ctx, query, t.ID, existingTodo.Title, existingTodo.Completed, existingTodo.Order).Scan()
	if err != nil {
		return err
	}

	t.Url = t.URL()

	return nil
}
