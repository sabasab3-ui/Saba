package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sabasab3-ui/saba/internal/intelligence"
)

type analyzeRequest struct {
	Question string                `json:"question"`
	Sources  []intelligence.Source `json:"sources"`
}

type analyzeResponse struct {
	Question    string   `json:"question"`
	Findings    []string `json:"findings"`
	Decision    string   `json:"decision"`
	Confidence  float64  `json:"confidence"`
	SourceCount int      `json:"source_count"`
}

func main() {
	analyzer := intelligence.SmartAnalyzer{}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"service": "SABA Intelligence",
		})
	})

	http.HandleFunc("/analyze", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST required", http.StatusMethodNotAllowed)
			return
		}

		var req analyzeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		findings, decision, confidence, err :=
			analyzer.Analyze(r.Context(), req.Question, req.Sources)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(analyzeResponse{
			Question:    req.Question,
			Findings:    findings,
			Decision:    decision,
			Confidence:  confidence,
			SourceCount: len(req.Sources),
		})
	})

	log.Println("SABA Intelligence running on http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
