package intelligence

import (
	"context"
	"errors"
	"strings"
)

// Source is a piece of evidence collected during an intelligence run.
type Source struct {
	Title   string
	URL     string
	Content string
}

// Report is the structured output of SABA Intelligence.
type Report struct {
	Question   string
	Summary    string
	Findings   []string
	Decision   string
	Confidence float64
	Sources    []Source
}

// Researcher collects evidence for an intelligence question.
type Researcher interface {
	Research(ctx context.Context, question string) ([]Source, error)
}

// Analyzer converts research evidence into findings and a decision.
type Analyzer interface {
	Analyze(
		ctx context.Context,
		question string,
		sources []Source,
	) (
		findings []string,
		decision string,
		confidence float64,
		err error,
	)
}

// Engine coordinates:
// Question -> Research -> Analysis -> Decision -> Report.
type Engine struct {
	Researcher Researcher
	Analyzer   Analyzer
}

func NewEngine(researcher Researcher, analyzer Analyzer) *Engine {
	return &Engine{
		Researcher: researcher,
		Analyzer:   analyzer,
	}
}

func (e *Engine) Run(
	ctx context.Context,
	question string,
) (*Report, error) {

	question = strings.TrimSpace(question)

	if question == "" {
		return nil, ErrEmptyQuestion
	}

	if e.Researcher == nil {
		return nil, errors.New(
			"intelligence researcher is not configured",
		)
	}

	if e.Analyzer == nil {
		return nil, errors.New(
			"intelligence analyzer is not configured",
		)
	}

	sources, err := e.Researcher.Research(
		ctx,
		question,
	)
	if err != nil {
		return nil, err
	}

	findings, decision, confidence, err :=
		e.Analyzer.Analyze(
			ctx,
			question,
			sources,
		)

	if err != nil {
		return nil, err
	}

	return &Report{
		Question:   question,
		Summary:    buildSummary(findings),
		Findings:   findings,
		Decision:   decision,
		Confidence: clampConfidence(confidence),
		Sources:    sources,
	}, nil
}

func buildSummary(findings []string) string {
	if len(findings) == 0 {
		return "No findings were produced."
	}

	return strings.Join(findings, " ")
}

func clampConfidence(value float64) float64 {
	if value < 0 {
		return 0
	}

	if value > 1 {
		return 1
	}

	return value
}

var ErrEmptyQuestion = errors.New(
	"intelligence question cannot be empty",
)
