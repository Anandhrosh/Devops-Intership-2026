package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
)

func main() {
    // Docker health check
    if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
        resp, err := http.Get("http://127.0.0.1:3000/health")
        if err != nil {
            os.Exit(1)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            os.Exit(1)
        }
        os.Exit(0)
    }

    // Main application page
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Welcome to Sprint 1 Docker Project!")
    })

    // Health endpoint
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprintln(w, "OK")
    })

    log.Println("Server running on port 3000")
    log.Fatal(http.ListenAndServe(":3000", nil))
}