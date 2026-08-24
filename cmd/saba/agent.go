package main

import (
	"encoding/json"
	"github.com/sabasab3-ui/saba/internal/gateway"
	"net/http"
)

func init() {
	g := gateway.New()
	http.HandleFunc("/agent", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST required", http.StatusMethodNotAllowed)
			return
		}
		var req gateway.AgentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(g.Route(req))
	})
}
