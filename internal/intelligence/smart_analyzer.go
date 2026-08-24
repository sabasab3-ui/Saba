package intelligence

import "context"

// SmartAnalyzer connects SABA's evidence analysis with its reasoning layer.
type SmartAnalyzer struct{}

// Analyze performs structured analysis and evidence-based reasoning.
func (a SmartAnalyzer) Analyze(
	ctx context.Context,
	question string,
	sources []Source,
) ([]string, string, float64, error) {
	if err := ctx.Err(); err != nil {
		return nil, "", 0, err
	}

	analysis := Analyze(sources)
	reasoning := Reason(sources)

	confidence := analysis.Confidence
	if reasoning.Agreement > 0 {
		confidence = (confidence + reasoning.Agreement) / 2
	}

	if confidence > 1 {
		confidence = 1
	}

	return analysis.KeyPoints, reasoning.Decision, confidence, nil
}
