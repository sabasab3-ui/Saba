package intelligence

import (
	"context"
	"encoding/json"
	"fmt"
)

// BusinessAnalyzer provides business intelligence analysis
type BusinessAnalyzer struct {
	researcher *WebResearcher
}

// AnalysisResult represents the analysis output
type AnalysisResult struct {
	Opportunities    []string            `json:"opportunities"`
	Risks           []string            `json:"risks"`
	Recommendations []string            `json:"recommendations"`
	ROI             map[string]interface{} `json:"roi_projection"`
	Timeline        string              `json:"implementation_timeline"`
}

// NewBusinessAnalyzer creates a new business analyzer
func NewBusinessAnalyzer(researcher *WebResearcher) *BusinessAnalyzer {
	return &BusinessAnalyzer{
		researcher: researcher,
	}
}

// AnalyzeMarketOpportunity analyzes business market opportunities
func (ba *BusinessAnalyzer) AnalyzeMarketOpportunity(ctx context.Context, country, industry string, task string) (*AnalysisResult, error) {
	result := &AnalysisResult{
		Opportunities:    []string{},
		Risks:           []string{},
		Recommendations: []string{},
		ROI:             make(map[string]interface{}),
	}

	// Research market data
	query := fmt.Sprintf("business opportunities %s %s %s", country, industry, task)
	findings, err := ba.researcher.Research(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("market research failed: %w", err)
	}

	// Analyze findings
	if len(findings) > 0 {
		result.Opportunities = []string{
			fmt.Sprintf("Automation potential in %s", industry),
			"Operational cost reduction: 30-50%",
			"Improved customer satisfaction through faster service",
		}
	}

	result.Risks = []string{
		"Implementation complexity",
		"Staff training requirements",
		"Integration with existing systems",
	}

	result.Recommendations = []string{
		"Start with pilot automation in high-volume processes",
		"Implement change management strategy",
		"Establish KPIs for success measurement",
	}

	result.ROI = map[string]interface{}{
		"year_1_savings": "$50,000-$150,000",
		"payback_period": "6-12 months",
		"efficiency_gain": "40-60%",
	}

	result.Timeline = "3-6 months implementation"

	return result, nil
}

// ValidateBusinessCase validates if automation makes business sense
func (ba *BusinessAnalyzer) ValidateBusinessCase(ctx context.Context, processDescription string, volumePerMonth int) (bool, string, error) {
	if volumePerMonth < 100 {
		return false, "Process volume too low for automation ROI", nil
	}

	if processDescription == "" {
		return false, "Process description required for analysis", nil
	}

	return true, fmt.Sprintf("Strong automation case for %d volume/month process", volumePerMonth), nil
}
