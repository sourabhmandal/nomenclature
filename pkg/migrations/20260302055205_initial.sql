-- Create "translations" table
CREATE TABLE "translations" (
  "id" serial NOT NULL,
  "normalized_hash" character(64) NOT NULL,
  "language_code" character varying(10) NOT NULL,
  "original_text" text NOT NULL,
  "translated_text" text NOT NULL,
  "confidence_score" numeric(4,3) NULL,
  "provider" character varying(50) NULL,
  "created_at" timestamp NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "translations_normalized_hash_key" UNIQUE ("normalized_hash"),
  CONSTRAINT "translations_normalized_hash_language_code_key" UNIQUE ("normalized_hash", "language_code")
);
-- Create index "idx_translations_lookup" to table: "translations"
CREATE INDEX "idx_translations_lookup" ON "translations" ("normalized_hash", "language_code");
-- Create "users" table
CREATE TABLE "users" (
  "id" bigserial NOT NULL,
  "name" text NOT NULL,
  "email" text NOT NULL,
  "bio" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "users_email_key" UNIQUE ("email")
);
