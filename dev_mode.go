//go:build dev
// +build dev

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func init() {
	ctx := context.Background()

	//localstackContainer, err := localstack.Run(ctx,
	//	"localstack/localstack:3.1",
	//)

	//localstackContainer.Endpoint()

	c, err := postgres.Run(ctx,
		"postgres:15.3-alpine",
		postgres.WithInitScripts(filepath.Join(".", "testdata", "schema.sql")),
		postgres.WithDatabase("todos"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		panic(err)
	}

	connStr, err := c.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}

	// check the connection to the database
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	App.UsersConnection = connStr
	log.Println("Users database started successfully")

	createSampleTodos(connStr)

	// register a graceful shutdown to stop the dependencies when the application is stopped
	// only in development mode
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		// also use the shutdown function when the SIGTERM or SIGINT signals are received
		sig := <-gracefulStop
		fmt.Printf("caught sig: %+v\n", sig)
		err := shutdownDependencies(c)
		if err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}()
}

func createSampleTodos(connStr string) {
	todoRepository, err := NewTodosRepository(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Cannot create a Todos repository. Exiting")
	}

	var title = "I need your root, your RAM, and your CPU cycles"
	var completed bool = false
	var order int = 1

	todoRepository.Create(context.Background(), &Todo{
		ID:        uuid.NewString(),
		Title:     title,
		Completed: &completed,
		Order:     &order,
	})
}

// helper function to stop the dependencies
func shutdownDependencies(containers ...testcontainers.Container) error {
	ctx := context.Background()
	for _, c := range containers {
		err := c.Terminate(ctx)
		if err != nil {
			log.Println("Error terminating the backend dependency:", err)
			return err
		}
	}

	return nil
}
