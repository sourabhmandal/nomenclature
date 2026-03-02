CREATE TABLE companies (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE translations (
    id SERIAL PRIMARY KEY,
    company_id BIGINT REFERENCES companies(id) ON DELETE CASCADE,
    normalized_hash CHAR(64) NOT NULL UNIQUE,
    source_language VARCHAR(10) NOT NULL,
    target_language VARCHAR(10) NOT NULL,
    original_text TEXT NOT NULL,
    translated_text TEXT NOT NULL,
    confidence_score NUMERIC(4,3),
    provider VARCHAR(50),
    created_at TIMESTAMP DEFAULT now(),
    UNIQUE(normalized_hash, source_language, target_language, company_id)
);

CREATE INDEX idx_translations_lookup 
ON translations(normalized_hash, source_language, target_language, company_id);

CREATE TABLE users (
    id   BIGSERIAL PRIMARY KEY,
    name text      NOT NULL,
    email text UNIQUE NOT NULL,
    bio  text
);