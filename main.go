package main

import (
 "encoding/json"
 "fmt"
 "log"
 "net/http"
 "os"
 "sync"
)

var (
 healthStatus = struct {
  status string
  sync.Mutex
 }{status: "up"}
)

func main() {
 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "Привет, мир")
 })

 http.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  fmt.Fprintln(w, r.URL.String())
 })

 http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
  healthStatus.Lock()
  defer healthStatus.Unlock()
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(map[string]string{"status": healthStatus.status})
 })

 http.HandleFunc("/healthz/update", func(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
   http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
   return
  }
  var req map[string]string
  if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
   http.Error(w, "Invalid JSON", http.StatusBadRequest)
   return
  }
  healthStatus.Lock()
  defer healthStatus.Unlock()
  healthStatus.status = req["status"]
  w.WriteHeader(http.StatusOK)
 })

 port := os.Getenv("APP_PORT")
 if port == "" {
  port = "8080"
 }
 log.Println("Starting server on port", port)
 log.Fatal(http.ListenAndServe(":"+port, nil))
}
