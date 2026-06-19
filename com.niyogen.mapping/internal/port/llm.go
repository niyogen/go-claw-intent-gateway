package port

import (
	"context"

	"com.niyogen/openclaw/internal/domain"
)

// IntentDetector calls the LLM to classify user text into a registered intent.
type IntentDetector interface {
	DetectIntent(ctx context.Context, msg domain.UnifiedMessage, availableIntents []domain.IntentDef) (*IntentDetectionResult, error)
}

// IntentDetectionResult is the structured output from the LLM.
type IntentDetectionResult struct {
	IntentName string
	Params     map[string]interface{}
}

// FeatureAnalyzer handles Capability 5: Review-Driven Experimentation.
type FeatureAnalyzer interface {
	AnalyzeReviews(ctx context.Context, reviews []string) (*ExperimentProposal, error)
}

type ExperimentProposal struct {
	TopFeature          string
	Confidence          float64
	MentionCount        int
	SupportingQuotes    []string
	ExperimentObjective string
	Hypothesis          string
}
