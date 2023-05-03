ALTER TABLE "users"
  DROP CONSTRAINT fk_leaves_users;

ALTER TABLE "leaves"
  DROP user_id;
