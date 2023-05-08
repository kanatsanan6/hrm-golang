CREATE TABLE IF NOT EXISTS "companies" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

-- Create the trigger to execute the update_updated_at function on UPDATE
CREATE TRIGGER companies_updated_at_trigger
BEFORE UPDATE ON "companies"
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();
