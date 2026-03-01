package translation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// GoogleTranslateProvider implements TranslationProvider for Google Cloud Translate.
type GoogleTranslateProvider struct {
	APIKey    string
	ProjectID string
}

func NewGoogleTranslateProvider(googleTranslateAPIKey string, projectID string) *GoogleTranslateProvider {
	return &GoogleTranslateProvider{
		APIKey:    googleTranslateAPIKey,
		ProjectID: projectID, // set this in ENV
	}
}

func (g *GoogleTranslateProvider) Translate(req TranslationRequest) (TranslationResult, error) {
	endpoint := fmt.Sprintf("https://translation.googleapis.com/v3beta1/projects/%s/locations/global:translateText?key=%s", g.ProjectID, g.APIKey)
	payload := map[string]interface{}{
		"contents":           []string{req.Text},
		"sourceLanguageCode": req.SourceLang,
		"targetLanguageCode": req.TargetLang,
	}
	body, _ := json.Marshal(payload)
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return TranslationResult{}, err
	}
	defer resp.Body.Close()
	var result struct {
		Translations []struct {
			TranslatedText string `json:"translatedText"`
		} `json:"translations"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return TranslationResult{}, err
	}
	if len(result.Translations) == 0 {
		return TranslationResult{}, fmt.Errorf("no translation returned")
	}
	tr := result.Translations[0]
	return TranslationResult{
		Text:     tr.TranslatedText,
		Provider: "google",
		Original: req.Text,
	}, nil
}
