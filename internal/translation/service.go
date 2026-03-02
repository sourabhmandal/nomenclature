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

	hash := GenerateHash(req.CompanyID, req.SourceLanguage, req.TargetLanguage, req.Text)
	// Translation Memory Lookup
	savedTranslationResponse, err := s.translationRepository.GetTranslationByHash(ctx, &repository.GetTranslationByHashParams{
		CompanyID:      &req.CompanyID,
		NormalizedHash: hash,
	})
	if err == nil && savedTranslationResponse.ID != 0 {
		log.Printf("CACHE HIT: %s (hash) :: %s (text)", hash, savedTranslationResponse.OriginalText)
		return TranslationResult{
			CompanyID:      req.CompanyID,
			NormalizedHash: hash,
			SourceLanguage: savedTranslationResponse.SourceLanguage,
			TargetLanguage: savedTranslationResponse.TargetLanguage,
			Original:       savedTranslationResponse.OriginalText,
			Translated:     savedTranslationResponse.TranslatedText,
			Confidence:     float64(savedTranslationResponse.ConfidenceScore.Int.Int64()),
			Provider:       savedTranslationResponse.Provider,
		}, nil
	}

	// If not found in memory, call the translation provider
	translationResp, err := s.translator.Translate(translationproviders.TranslationRequest{
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
	savedTranslationResponse, err = s.translationRepository.SaveTranslationByHash(ctx, &repository.SaveTranslationByHashParams{
		CompanyID:      &req.CompanyID,
		NormalizedHash: hash,
		SourceLanguage: req.SourceLanguage,
		TargetLanguage: req.TargetLanguage,
		OriginalText:   translationResp.Original,
		TranslatedText: translationResp.Translated,
		ConfidenceScore: pgtype.Numeric{
			Int:   big.NewInt(int64(translationResp.Confidence)),
			Exp:   -2,
			Valid: true,
		},
		Provider: &translationResp.Provider,
	})
	if err != nil {
		log.Printf("Failed to save translation: %v", err)
		// Even if saving fails, we can return the translation result
		// but without the confidence score and provider information.
		return TranslationResult{
			CompanyID:      req.CompanyID,
			NormalizedHash: hash,
			SourceLanguage: savedTranslationResponse.SourceLanguage,
			TargetLanguage: savedTranslationResponse.TargetLanguage,
			Original:       req.Text,
			Translated:     translationResp.Translated,
			Confidence:     float64(savedTranslationResponse.ConfidenceScore.Int.Int64()),
			Provider:       savedTranslationResponse.Provider,
		}, nil
	}
	return TranslationResult{
		CompanyID:      req.CompanyID,
		NormalizedHash: hash,
		SourceLanguage: savedTranslationResponse.SourceLanguage,
		TargetLanguage: savedTranslationResponse.TargetLanguage,
		Original:       req.Text,
		Translated:     savedTranslationResponse.TranslatedText,
		Confidence:     float64(savedTranslationResponse.ConfidenceScore.Int.Int64()),
		Provider:       savedTranslationResponse.Provider,
	}, nil
}
