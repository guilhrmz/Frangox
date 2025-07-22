package main

import (
  "frangox/routes"
  "log"
  "net/http"
)

func main() {
  r := routes.SetupRoutes()
  log.Println("Servidor rodando em http://localhost:8080")
  http.ListenAndServe(":8080", r)
}
