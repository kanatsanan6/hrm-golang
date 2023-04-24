CREATE TABLE IF NOT EXISTS "users" (
  "id" bigserial PRIMARY KEY,
  "first_name" varchar(255) NOT NULL,
  "last_name" varchar(255) NOT NULL,
  "email" varchar(255) NOT NULL,
  "encrypted_password" varchar(255) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL
);

CREATE UNIQUE INDEX ON "users" ("email");
