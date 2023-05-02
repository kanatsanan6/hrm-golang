CREATE TABLE IF NOT EXISTS "leaves" (
  "id" bigserial PRIMARY KEY,
  "description" varchar(255) NOT NULL,
  "status" varchar NOT NULL DEFAULT 'pending' CHECK ("status" in ('pending', 'approved', 'rejected')),
  "start_date" timestamptz NOT NULL,
  "end_date" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL
);
