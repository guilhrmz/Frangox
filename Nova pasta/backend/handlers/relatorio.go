package handlers

import (
  "encoding/json"
  "fmt"
  "net/http"
  "os"
)

func RelatorioHandler(w http.ResponseWriter, r *http.Request) {
  total := 270.00
  pedidos := 6
  lucro := total * 0.3

  json.NewEncoder(w).Encode(map[string]interface{}{
    "total_vendas": total,
    "total_pedidos": pedidos,
    "lucro": lucro,
  })
}

func DownloadRelatorioHandler(w http.ResponseWriter, r *http.Request) {
  tipo := r.URL.Query().Get("tipo")

  var file string
  switch tipo {
  case "txt":
    file = gerarTXT()
  default:
    http.Error(w, "Formato não suportado", http.StatusBadRequest)
    return
  }

  http.ServeFile(w, r, file)
}

func gerarTXT() string {
  path := "relatorio.txt"
  f, _ := os.Create(path)
  defer f.Close()

  f.WriteString("RELATÓRIO DE VENDAS\n")
  f.WriteString("=====================\n")
  f.WriteString("Total de pedidos vendidos: 6\n")
  f.WriteString("Total faturado: R$270.00\n")
  f.WriteString("Lucro estimado: R$81.00\n")

  return path
}
