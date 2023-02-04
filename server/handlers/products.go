package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type Product struct {
	Id         int     `json:"id,omitempty"`
	Name       string  `json:"name,omitempty"`
	Price      float32 `json:"price,omitempty"`
	Image_Path string  `json:"image_path,omitempty"`
	SKU        string  `json:"sku,omitempty"`
}

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(db_conn *pgx.Conn) func(http.ResponseWriter, *http.Request) {
	if db_conn == nil {
		panic("Nil db_conn in GetProducts!")
	}
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Max-Age", "86400")
		rows, err := db_conn.Query(context.Background(), "SELECT * FROM frontend")
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
}
