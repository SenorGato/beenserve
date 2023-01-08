package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

func main() {
	// sm := mux.NewRouter()
	// urlExample := "postgres://tealacarte:smoke@localhost:5432/tealacarte"
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var name string
	var weight int64
	err = conn.QueryRow(context.Background(), "select name, weight from widgets where id=$1", 42).Scan(&name, &weight)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(name, weight)
}
