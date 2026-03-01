package translation

// TranslationRequest represents a request for translation.
type TranslationRequest struct {
	ProjectID  int64
	Text       string
	SourceLang string
	TargetLang string
}

// TranslationResult represents the result of a translation.
type TranslationResult struct {
	Text       string
	Confidence float64
	Provider   string
	Fallback   bool
	Original   string
}

// ConfidenceScore holds scoring details for a translation.
type ConfidenceScore struct {
	ProviderConfidence float64
	LanguageMatch      float64
	LengthRatio        float64
	SemanticSimilarity float64
	TotalScore         float64
}
