CREATE TABLE IF NOT EXISTS "companies" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL
);
