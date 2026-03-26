package translationproviders

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"google.golang.org/genai"
)

// GoogleTranslateProvider implements TranslationProvider for Google Cloud Translate.
type geminiProvider struct {
	gemini *genai.Client
}

func NewGoogleTranslateProvider(apiKey string) TranslationProvider {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatalf("Failed to create Gemini client: %v", err)
	}

	return &geminiProvider{
		gemini: client,
	}
}

func (g *geminiProvider) Translate(req TranslationRequest) (TranslationOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Minute)
	defer cancel()

	// Define strict response schema
	schema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"translations": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"original": {
							Type:        genai.TypeString,
							Description: "The original text exactly as provided",
						},
						"translated": {
							Type:        genai.TypeString,
							Description: "Translated text, or original if type is not string",
						},
						"type": {
							Type:        genai.TypeString,
							Description: "The type of the value: string, number, boolean, or link",
							Enum:        []string{"string", "number", "boolean", "link"},
						},
					},
					Required: []string{"original", "translated", "type"},
				},
			},
		},
		Required: []string{"translations"},
	}

	result, err := g.gemini.Models.GenerateContent(
		ctx,
		"gemini-flash-latest",
		genai.Text(buildPrompt(req)),
		&genai.GenerateContentConfig{
			ResponseMIMEType: "application/json",
			ResponseSchema:   schema,
			Temperature:      genai.Ptr[float32](0.1),
		},
	)

	if err != nil {
		return TranslationOutput{}, err
	}

	if err != nil {
		return TranslationOutput{}, fmt.Errorf("gemini API error: %w", err)
	}

	// Extract text from response
	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return TranslationOutput{}, fmt.Errorf("empty response from Gemini")
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return TranslationOutput{}, fmt.Errorf("empty response from Gemini")
	}

	rawJSON := result.Candidates[0].Content.Parts[0].Text

	var output TranslationOutput
	if err := json.Unmarshal([]byte(rawJSON), &output); err != nil {
		return TranslationOutput{}, fmt.Errorf("failed to parse response: %w\nRaw: %s", err, rawJSON)
	}

	// --- Print structured output ---
	structured, _ := json.MarshalIndent(output, "", "  ")
	fmt.Println("=== Structured Output ===")
	fmt.Println(string(structured))

	return output, nil
}

func buildPrompt(req TranslationRequest) string {
	// Serialize context segments to JSON for the prompt
	contextJSON, _ := json.MarshalIndent(req.TextWithContext, "", "  ")

	return fmt.Sprintf(`You are a website translation engine.

You will be given an array of context objects. Each object maps a text value to its TYPE (e.g. "string", "number", "link", "boolean").

RULES:
- Translate ONLY values of type "string"
- Keep values of type "number", "boolean", "link" exactly as-is (do not translate)
- Preserve original casing and punctuation style
- Return ONLY a valid JSON object with NO markdown, NO explanation
- Format: {"original text": "translated text", ...}
- Every key from every segment object must appear in the output

Source language: %s
Target language: %s

Context segments (array of {text: type} maps):
%s

Respond with only the JSON object.`,
		req.SourceLanguage,
		req.TargetLanguage,
		string(contextJSON),
	)
}

// Helper: convert structured output → final {"original": "translated"} map
func toFinalMap(output *TranslationOutput) map[string]string {
	result := make(map[string]string, len(output.Translations))
	for _, pair := range output.Translations {
		result[pair.Original] = pair.Translated
	}
	return result
}
