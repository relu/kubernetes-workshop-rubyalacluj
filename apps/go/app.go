package main

import (
  "net/http"
  "os"
  "runtime"
  "path/filepath"
  "fmt"
  "time"
)

const PORT = 3000

func main() {
  _, name, _, _ := runtime.Caller(0)
  name = filepath.Base(name)
  if value, ok := os.LookupEnv("NAME"); ok {
    name = value
  }

  mux := http.NewServeMux()
  mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("%s [%s] \"%s %s HTTP/1.1\" %d %s\n", r.RemoteAddr, time.Now().Format("2006-01-02T15:04:05-0700"), r.Method, r.URL.Path, 200, r.UserAgent())
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(fmt.Sprintf("Hello from %s", name)))
  })
  fmt.Printf("listening on %d\n", PORT)
  http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", PORT), mux)
}
