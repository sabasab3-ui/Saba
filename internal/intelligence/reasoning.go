package intelligence

import "strings"

// ReasoningResult contains SABA's evidence-based reasoning.
type ReasoningResult struct {
	Agreement   float64
	Uncertainty float64
	Decision    string
	Explanation string
}

// Reason evaluates evidence agreement and uncertainty.
func Reason(sources []Source) ReasoningResult {
	if len(sources) == 0 {
		return ReasoningResult{
			Uncertainty: 1,
			Decision:    "Insufficient evidence",
			Explanation: "SABA has no evidence to evaluate.",
		}
	}

	seen := make(map[string]int)
	for _, source := range sources {
		text := strings.ToLower(strings.TrimSpace(source.Title + " " + source.Content))
		if text == "" {
			continue
		}

		words := strings.Fields(text)
		for _, word := range words {
			if len(word) >= 5 {
				seen[word]++
			}
		}
	}

	agreement := 0.0
	if len(sources) >= 2 {
		agreement = 0.5
	}
	if len(sources) >= 4 {
		agreement = 0.75
	}

	for _, count := range seen {
		if count >= 2 {
			agreement += 0.05
			break
		}
	}

	if agreement > 1 {
		agreement = 1
	}

	uncertainty := 1 - agreement

	decision := "Evidence is insufficient for a strong conclusion."
	if agreement >= 0.75 {
		decision = "Evidence shows meaningful agreement."
	} else if agreement >= 0.5 {
		decision = "Evidence provides moderate support."
	}

	return ReasoningResult{
		Agreement:   agreement,
		Uncertainty: uncertainty,
		Decision:    decision,
		Explanation: "SABA compared the available evidence and estimated agreement and uncertainty.",
	}
}
