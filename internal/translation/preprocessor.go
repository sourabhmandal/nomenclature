package translation

import (
	"regexp"
	"strings"
)

var (
	linkPattern       = regexp.MustCompile(`(?i)^(https?://|www\.)[^\s]+$`)
	booleanPattern    = regexp.MustCompile(`(?i)^(true|false|yes|no|on|off|y|n)$`)
	numberPattern     = regexp.MustCompile(`^[+-]?(?:\d+\.?\d*|\.\d+)$`)
	emailPattern      = regexp.MustCompile(`(?i)^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	uuidPattern       = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	datePattern       = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	alphaPattern      = regexp.MustCompile(`^[A-Za-z]+$`)
	alphanumericToken = regexp.MustCompile(`^[A-Za-z0-9]+$`)
)

// DetectInputType splits input by spaces, classifies each token using regex,
// and returns a map of token -> detected type.
func DetectInputType(input string) map[string]string {
	tokens := strings.Fields(input)
	result := make(map[string]string, len(tokens))
	if len(tokens) == 0 {
		return result
	}

	for _, token := range tokens {
		normalized := strings.TrimSpace(strings.Trim(token, `.,;:!?()[]{}"'`))
		if normalized == "" {
			continue
		}
		result[normalized] = detectTokenType(normalized)
	}

	return result
}

func detectTokenType(token string) string {
	normalized := strings.TrimSpace(strings.Trim(token, `.,;:!?()[]{}"'`))
	if normalized == "" {
		return "string"
	}

	switch {
	case linkPattern.MatchString(normalized):
		return "link"
	case booleanPattern.MatchString(normalized):
		return "boolean"
	case numberPattern.MatchString(normalized):
		return "number"
	case emailPattern.MatchString(normalized):
		return "email"
	case uuidPattern.MatchString(normalized):
		return "uuid"
	case datePattern.MatchString(normalized):
		return "date"
	case alphaPattern.MatchString(normalized):
		return "string"
	case alphanumericToken.MatchString(normalized):
		return "string"
	default:
		return "string"
	}
}
