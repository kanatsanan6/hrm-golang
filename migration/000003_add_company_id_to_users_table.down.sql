ALTER TABLE "users"
  DROP CONSTRAINT fk_users_companies;

ALTER TABLE "users"
  DROP company_id;
