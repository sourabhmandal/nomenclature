package translation

import (
	"errors"
	"log"
)

// TranslationService implements the fail-safe translation pipeline.
type TranslationService struct {
	ProviderPrimary TranslationProvider
	Memory          TranslationMemory
}

func NewTranslationService(provider TranslationProvider, memory TranslationMemory) *TranslationService {
	return &TranslationService{
		ProviderPrimary: provider,
		Memory:          memory,
	}
}

func (s *TranslationService) Translate(req TranslationRequest) (TranslationResult, error) {
	// Translation Memory Lookup
	if res, found := s.Memory.Lookup(req); found {
		return res, nil
	}

	// Primary Provider
	res, err := s.ProviderPrimary.Translate(req)
	if err != nil {
		log.Printf("Primary provider failed: %v", err)
		return TranslationResult{Text: req.Text, Original: req.Text}, errors.New("translation failed")
	}

	// Store in Memory
	s.Memory.Store(req, res)
	return res, nil
}
