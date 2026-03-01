package translation

// TranslationProvider defines the interface for translation providers.
type TranslationProvider interface {
	Translate(req TranslationRequest) (TranslationResult, error)
}

// TranslationMemory defines the interface for translation memory lookup.
type TranslationMemory interface {
	Lookup(req TranslationRequest) (TranslationResult, bool)
	Store(req TranslationRequest, res TranslationResult) error
}
