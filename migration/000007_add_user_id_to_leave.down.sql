ALTER TABLE "users"
  DROP CONSTRAINT IF EXISTS fk_leaves_users;

ALTER TABLE "leaves"
  DROP user_id;
