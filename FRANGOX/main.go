package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql" // ou _ "github.com/lib/pq" para PostgreSQL
)

type Pedido struct {
	ID          int
	Cliente     string
	Produto     string
	Entrega     string
	Preco       float64
	Status      string
	StatusClass string
}

var db *sql.DB

func main() {
	var err error
	// Conexão MySQL (substitua pelas suas credenciais)
	db, err = sql.Open("mysql", "root:senha@tcp(127.0.0.1:3306)/frangox")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", home)

	log.Println("Servidor rodando em http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT id, cliente, produto, entrega, preco, status FROM pedidos ORDER BY id DESC`)
	if err != nil {
		http.Error(w, "Erro ao buscar pedidos", 500)
		return
	}
	defer rows.Close()

	var pedidos []Pedido

	for rows.Next() {
		var p Pedido
		err := rows.Scan(&p.ID, &p.Cliente, &p.Produto, &p.Entrega, &p.Preco, &p.Status)
		if err != nil {
			continue
		}
		if p.Status == "VENDIDO" {
			p.StatusClass = "VENDIDO"
		} else {
			p.StatusClass = "PENDENTE"
		}
		pedidos = append(pedidos, p)
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, struct{ Pedidos []Pedido }{pedidos})
}
