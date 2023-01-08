package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Product struct {
	Id        int                `json:"id,omitempty"`
	Name      string             `json:"name,omitempty"`
	Price     float32            `json:"price,omitempty"`
	SKU       string             `json:"sku,omitempty"`
	CreatedOn pgtype.Timestamptz `json:"created_on,omitempty"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at,omitempty"`
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
	products, err := pgx.CollectRows(rows, pgx.RowToStructByPos[Product])
	if err != nil {
		fmt.Printf("Collect rows error: %v", err)
		return
	}

	d, err := json.Marshal(products)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(d)
}
