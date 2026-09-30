// Command demoserver supplies invented API responses for an account-free VHS recording.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	addr := "127.0.0.1:8644"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}
	server := &http.Server{
		Addr: addr, Handler: demoHandler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	fmt.Fprintln(os.Stderr, "demo API listening on", addr)
	if err := server.ListenAndServe(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(value)
}

type rewrite struct {
	Domain string `json:"domain"`
	Answer string `json:"answer"`
}

func demoHandler() http.Handler {
	var mu sync.Mutex
	rewrites := []rewrite{{Domain: "dashboard.example.test", Answer: "192.0.2.10"}, {Domain: "storage.example.test", Answer: "192.0.2.20"}}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		switch r.Method + " " + r.URL.Path {
		case "GET /control/status":
			writeJSON(w, map[string]any{"version": "v0.107.0", "running": true, "protection_enabled": true, "dns_port": 53, "http_port": 3000, "dns_addresses": []string{"192.0.2.2"}})
		case "GET /control/stats":
			writeJSON(w, map[string]any{"num_dns_queries": 12480, "num_blocked_filtering": 3120,
				"num_replaced_safebrowsing": 12, "num_replaced_parental": 0, "avg_processing_time": 0.0024,
				"top_queried_domains": []any{map[string]any{"docs.example.test": 240}}, "time_units": "hours"})
		case "GET /control/rewrite/list":
			writeJSON(w, rewrites)
		case "POST /control/rewrite/add":
			var entry rewrite
			if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			// Keep writes in memory so the final list visibly includes the requested rewrite.
			rewrites = append(rewrites, entry)
			writeJSON(w, entry)
		default:
			http.NotFound(w, r)
		}
	})
}
