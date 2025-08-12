package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type OrderItem struct {
	ProductID int     `json:"productId"`
	Name      string  `json:"name"`
	UnitPrice float64 `json:"unitPrice"`
	Quantity  int     `json:"quantity"`
}

type Order struct {
	ID        int         `json:"id"`
	Customer  string      `json:"customer"`
	Items     []OrderItem `json:"items"`
	Total     float64     `json:"total"`
	Status    string      `json:"status"` // "aberto" | "finalizado"
	CreatedAt time.Time   `json:"createdAt"`
}

var (
	products = []Product{
		{ID: 1, Name: "Frango Assado", Price: 45.00}, // produto pré-definido
	}
	productsMu sync.Mutex

	orders   = make([]Order, 0)
	ordersMu sync.Mutex
	nextID   = 1
)

func main() {
	http.HandleFunc("/api/products", productsHandler)   // GET
	http.HandleFunc("/api/orders", ordersHandler)       // GET, POST
	http.HandleFunc("/api/orders/", orderActionHandler) // finalize: POST /api/orders/{id}/finalize

	// arquivos estáticos
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	log.Println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	productsMu.Lock()
	defer productsMu.Unlock()
	json.NewEncoder(w).Encode(products)
}

func ordersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		ordersMu.Lock()
		out := make([]Order, len(orders))
		copy(out, orders)
		ordersMu.Unlock()
		json.NewEncoder(w).Encode(out)
	case http.MethodPost:
		// Estrutura esperada: { customer: "", items: [{ productId: 1, quantity: 2 }, ...] }
		var req struct {
			Customer string `json:"customer"`
			Items    []struct {
				ProductID int `json:"productId"`
				Quantity  int `json:"quantity"`
			} `json:"items"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if len(req.Items) == 0 {
			http.Error(w, "no items", http.StatusBadRequest)
			return
		}

		// montar itens do pedido e calcular total no servidor
		productsMu.Lock()
		defer productsMu.Unlock()

		var orderItems []OrderItem
		var total float64
		for _, it := range req.Items {
			if it.Quantity <= 0 {
				continue
			}
			var found *Product
			for _, p := range products {
				if p.ID == it.ProductID {
					tmp := p
					found = &tmp
					break
				}
			}
			if found == nil {
				http.Error(w, "product not found: "+strconv.Itoa(it.ProductID), http.StatusBadRequest)
				return
			}
			orderItems = append(orderItems, OrderItem{
				ProductID: found.ID,
				Name:      found.Name,
				UnitPrice: found.Price,
				Quantity:  it.Quantity,
			})
			total += found.Price * float64(it.Quantity)
		}
		if len(orderItems) == 0 {
			http.Error(w, "no valid items", http.StatusBadRequest)
			return
		}

		ordersMu.Lock()
		order := Order{
			ID:        nextID,
			Customer:  req.Customer,
			Items:     orderItems,
			Total:     total,
			Status:    "aberto",
			CreatedAt: time.Now(),
		}
		nextID++
		// insere no topo (mais recentes primeiro)
		orders = append([]Order{order}, orders...)
		ordersMu.Unlock()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(order)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func orderActionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// espera: /api/orders/{id}/finalize
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/orders/")
	// path agora: e.g. "3/finalize"
	if strings.HasSuffix(path, "/finalize") {
		idStr := strings.TrimSuffix(path, "/finalize")
		idStr = strings.Trim(idStr, "/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		ordersMu.Lock()
		defer ordersMu.Unlock()
		for i := range orders {
			if orders[i].ID == id {
				orders[i].Status = "finalizado"
				json.NewEncoder(w).Encode(orders[i])
				return
			}
		}
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	http.NotFound(w, r)
}
