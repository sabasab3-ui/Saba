package gateway

import (
	"context"
	"fmt"

	"github.com/sabasab3-ui/saba/internal/intelligence"
)

type AgentRequest struct {
	Agent string `json:"agent"`
	Task  string `json:"task"`
	Input string `json:"input"`
}

type AgentResponse struct {
	Agent      string   `json:"agent"`
	Status     string   `json:"status"`
	Output     string   `json:"output"`
	Findings   []string `json:"findings,omitempty"`
	Decision   string   `json:"decision,omitempty"`
	Confidence float64  `json:"confidence,omitempty"`
}

type Gateway struct{}

func New() *Gateway {
	return &Gateway{}
}

func (g *Gateway) Route(req AgentRequest) AgentResponse {
	analyzer := intelligence.SmartAnalyzer{}

	findings, decision, confidence, err :=
		analyzer.Analyze(context.Background(), req.Input, nil)

	if err != nil {
		return AgentResponse{
			Agent:  req.Agent,
			Status: "error",
			Output: err.Error(),
		}
	}

	return AgentResponse{
		Agent:      req.Agent,
		Status:     "analyzed",
		Output:     fmt.Sprintf("%s agent analyzed task '%s'.", req.Agent, req.Task),
		Findings:   findings,
		Decision:   decision,
		Confidence: confidence,
	}
}
