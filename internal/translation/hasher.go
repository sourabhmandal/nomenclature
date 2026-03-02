package translation

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
)

func GenerateHash(companyID int64, languageCode string, text string) string {
	normalized := strings.TrimSpace(strings.ToLower(text))
	data := strconv.FormatInt(companyID, 10) + ":" + languageCode + ":" + normalized
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
