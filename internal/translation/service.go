package translation

import (
	"context"
	"errors"
	"log"
	"math/big"
	"nomenclature/internal/repository"
	"nomenclature/internal/translationproviders"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// TranslationService implements the fail-safe translation pipeline.
type TranslationService interface {
	Translate(req TranslationRequest) (TranslationResult, error)
}

type translationService struct {
	translationRepository repository.Querier
	translator            translationproviders.TranslationProvider
}

func NewTranslationService(provider translationproviders.TranslationProvider, translationRepository repository.Querier) TranslationService {
	return &translationService{
		translator:            provider,
		translationRepository: translationRepository,
	}
}

func (s *translationService) Translate(req TranslationRequest) (TranslationResult, error) {
	// context.Background() can be used here if needed for timeouts or cancellation in future enhancements.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	hash := GenerateHash(req.CompanyID, req.SourceLanguage, req.Text)
	// Translation Memory Lookup
	res, err := s.translationRepository.GetTranslationByHash(ctx, &repository.GetTranslationByHashParams{
		CompanyID:      &req.CompanyID,
		NormalizedHash: hash,
	})
	if err == nil {
		return TranslationResult{
			Original:   res.OriginalText,
			Translated: res.TranslatedText,
		}, nil
	}

	if res.ID != 0 {
		log.Printf("CACHE HIT: %s (hash) :: %s (text)", hash, res.OriginalText)
		return TranslationResult{
			Original:   res.OriginalText,
			Translated: res.TranslatedText,
		}, nil
	}

	// If not found in memory, call the translation provider
	translation, err := s.translator.Translate(translationproviders.TranslationRequest{
		CompanyID:      req.CompanyID,
		Text:           req.Text,
		SourceLanguage: req.SourceLanguage,
		TargetLanguage: req.TargetLanguage,
	})
	if err != nil {
		log.Printf("Primary provider failed: %v", err)
		return TranslationResult{
			Original:   req.Text,
			Translated: req.Text,
		}, errors.New("translation failed")
	}

	// Store in Memory
	s.translationRepository.SaveTranslationByHash(ctx, &repository.SaveTranslationByHashParams{
		CompanyID:       &req.CompanyID,
		NormalizedHash:  hash,
		SourceLanguage:  req.SourceLanguage,
		TargetLanguage:  req.TargetLanguage,
		OriginalText:    translation.Original,
		TranslatedText:  translation.Translated,
		ConfidenceScore: pgtype.Numeric{Valid: true, Int: big.NewInt(90)},
		Provider:        &translation.Provider,
	})
	return TranslationResult{
		Original:   req.Text,
		Translated: translation.Translated,
		Confidence: translation.Confidence,
		Provider:   translation.Provider,
	}, nil
}
