package intelligence

import (
	"sort"
	"strings"
)

// Analysis contains SABA's structured assessment of collected information.
type Analysis struct {
	Summary       string
	KeyPoints     []string
	Confidence    float64
	SourceCount   int
	Recommendations []string
}

// Analyze examines collected sources and produces a structured result.
func Analyze(sources []Source) Analysis {
	if len(sources) == 0 {
		return Analysis{
			Summary:    "No evidence was provided for analysis.",
			Confidence: 0,
			SourceCount: 0,
		}
	}

	keyPoints := make([]string, 0, len(sources))
	seen := make(map[string]bool)

	for _, source := range sources {
		title := strings.TrimSpace(source.Title)
		content := strings.TrimSpace(source.Content)

		if title == "" && content == "" {
			continue
		}

		point := title
		if point == "" {
			point = content
		}

		point = strings.Join(strings.Fields(point), " ")

		if point != "" && !seen[strings.ToLower(point)] {
			keyPoints = append(keyPoints, point)
			seen[strings.ToLower(point)] = true
		}
	}

	// Put the strongest-looking evidence first.
	sort.SliceStable(keyPoints, func(i, j int) bool {
		return len(keyPoints[i]) > len(keyPoints[j])
	})

	confidence := calculateConfidence(len(sources), len(keyPoints))

	summary := buildSummary(keyPoints)

	return Analysis{
		Summary:         summary,
		KeyPoints:       keyPoints,
		Confidence:      confidence,
		SourceCount:     len(sources),
		Recommendations: buildRecommendations(confidence),
	}
}

func calculateConfidence(sourceCount, usefulPoints int) float64 {
	if sourceCount == 0 {
		return 0
	}

	confidence := 0.35

	if sourceCount >= 2 {
		confidence += 0.20
	}

	if sourceCount >= 4 {
		confidence += 0.15
	}

	if usefulPoints >= 3 {
		confidence += 0.10
	}

	if usefulPoints >= 5 {
		confidence += 0.10
	}

	if confidence > 1 {
		confidence = 1
	}

	return confidence
}

func buildEvidenceSummary(points []string) string {
	if len(points) == 0 {
		return "SABA found no useful evidence to summarize."
	}

	if len(points) == 1 {
		return "SABA found one useful piece of evidence: " + points[0]
	}

	return "SABA analyzed multiple pieces of evidence and identified " +
		string(rune('0'+min(len(points), 9))) +
		" key information points."
}

func buildRecommendations(confidence float64) []string {
	if confidence < 0.5 {
		return []string{
			"Collect more sources before making an important decision.",
			"Verify important claims independently.",
		}
	}

	return []string{
		"Compare information across multiple sources.",
		"Verify important claims before acting.",
		"Use the strongest evidence to guide the next decision.",
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
