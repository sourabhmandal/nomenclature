-- Create "companies" table
CREATE TABLE "companies" (
  "id" serial NOT NULL,
  "name" character varying(255) NOT NULL,
  "created_at" timestamp NULL DEFAULT now(),
  PRIMARY KEY ("id")
);
-- Create "users" table
CREATE TABLE "users" (
  "id" bigserial NOT NULL,
  "name" text NOT NULL,
  "email" text NOT NULL,
  "bio" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "users_email_key" UNIQUE ("email")
);
-- Create "translations" table
CREATE TABLE "translations" (
  "id" serial NOT NULL,
  "company_id" bigint NULL,
  "normalized_hash" character(64) NOT NULL,
  "source_language" character varying(10) NOT NULL,
  "target_language" character varying(10) NOT NULL,
  "original_text" text NOT NULL,
  "translated_text" text NOT NULL,
  "confidence_score" numeric(4,3) NULL,
  "provider" character varying(50) NULL,
  "created_at" timestamp NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "translations_normalized_hash_key" UNIQUE ("normalized_hash"),
  CONSTRAINT "translations_normalized_hash_source_language_target_languag_key" UNIQUE ("normalized_hash", "source_language", "target_language", "company_id"),
  CONSTRAINT "translations_company_id_fkey" FOREIGN KEY ("company_id") REFERENCES "companies" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_translations_lookup" to table: "translations"
CREATE INDEX "idx_translations_lookup" ON "translations" ("normalized_hash", "source_language", "target_language", "company_id");
