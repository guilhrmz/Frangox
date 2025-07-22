package routes

import (
  "frangox/handlers"
  "github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
  r := mux.NewRouter()
  r.HandleFunc("/pedidos", handlers.ListarPedidos).Methods("GET")
  r.HandleFunc("/relatorio", handlers.RelatorioHandler).Methods("GET")
  r.HandleFunc("/relatorio/download", handlers.DownloadRelatorioHandler).Methods("GET")
  return r
}
