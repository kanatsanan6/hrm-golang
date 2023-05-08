CREATE TABLE IF NOT EXISTS "leaves" (
  "id" bigserial PRIMARY KEY,
  "description" varchar(255) NOT NULL,
  "status" varchar NOT NULL DEFAULT 'pending' CHECK ("status" in ('pending', 'approved', 'rejected')),
  "start_date" timestamptz NOT NULL,
  "end_date" timestamptz NOT NULL,
  "leave_type" varchar(255) NOT NULL CHECK ("leave_type" in ('vacation_leave', 'extra_vacation', 'sick_leave', 'business_leave')),
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

-- Create the trigger to execute the update_updated_at function on UPDATE
CREATE TRIGGER leaves_updated_at_trigger
BEFORE UPDATE ON "leaves"
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();
