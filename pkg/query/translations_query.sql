-- name: GetTranslationByHash :one
SELECT * FROM translations
WHERE company_id = $1 AND normalized_hash = $2 LIMIT 1;

-- name: SaveTranslationByHash :one
INSERT INTO translations (company_id, normalized_hash, source_language, target_language, original_text, translated_text, confidence_score, provider)
SELECT $1, $2, $3, $4, $5, $6, $7, $8
WHERE NOT EXISTS (
	SELECT 1 FROM translations WHERE company_id = $1 AND normalized_hash = $2 AND source_language = $3 AND target_language = $4
)
RETURNING *;
