package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/sabasab3-ui/saba/internal/agent"
	"github.com/sabasab3-ui/saba/internal/intelligence"
)

type analyzeRequest struct {
	Question string                `json:"question"`
	Sources  []intelligence.Source `json:"sources,omitempty"`
}

type analyzeResponse struct {
	Question    string                `json:"question"`
	Summary     string                `json:"summary"`
	Findings    []string              `json:"findings"`
	Decision    string                `json:"decision"`
	Confidence  float64               `json:"confidence"`
	SourceCount int                   `json:"source_count"`
	Sources     []intelligence.Source `json:"sources,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func main() {
	analyzer := intelligence.SmartAnalyzer{}
	researcher := intelligence.NewWebResearcher()
	engine := intelligence.NewEngine(researcher, analyzer)
	toolkit := agent.NewDefaultToolKit()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"service": "SABA Intelligence",
			"status":  "ok",
			"message": "SABA is online.",
			"endpoints": []string{
				"/health",
				"/tools",
				"/analyze",
			},
		})
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"status":  "ok",
			"service": "SABA Intelligence",
		})
	})

	// /tools exposes the unified tool package so the client/dashboard can
	// discover which capabilities are currently registered.
	http.HandleFunc("/tools", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
				"error": "GET required",
			})
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"count": len(toolkit.Catalog()),
			"tools": toolkit.Catalog(),
		})
	})

	http.HandleFunc("/analyze", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
				"error": "POST required",
			})
			return
		}

		var req analyzeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "invalid JSON",
			})
			return
		}

		if req.Question == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "question is required",
			})
			return
		}

		var (
			findings   []string
			decision   string
			confidence float64
			sources    []intelligence.Source
			summary    string
			err        error
		)

		if len(req.Sources) > 0 {
			sources = req.Sources
			findings, decision, confidence, err =
				analyzer.Analyze(r.Context(), req.Question, sources)
			summary = "SABA analyzed the supplied evidence."
		} else {
			report, runErr := engine.Run(r.Context(), req.Question)
			if runErr == nil {
				findings = report.Findings
				decision = report.Decision
				confidence = report.Confidence
				sources = report.Sources
				summary = report.Summary
			}
			err = runErr
		}

		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}

		writeJSON(w, http.StatusOK, analyzeResponse{
			Question:    req.Question,
			Summary:     summary,
			Findings:    findings,
			Decision:    decision,
			Confidence:  confidence,
			SourceCount: len(sources),
			Sources:     sources,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("SABA Intelligence running on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
