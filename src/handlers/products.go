package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Product struct {
	Id        int
	Name      string
	Price     float32
	SKU       string
	CreatedOn pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, h *http.Request) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "SELECT * FROM products")
	defer rows.Close()

	var rowSlice []Product
	for rows.Next() {
		var r Product
		err := rows.Scan(&r.Id, &r.Name, &r.Price, &r.SKU, &r.CreatedOn, &r.UpdatedAt)
		if err != nil {
			log.Fatal(err)
		}
		rowSlice = append(rowSlice, r)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(rowSlice)
}
