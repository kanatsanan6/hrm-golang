ALTER TABLE "users"
  ADD role varchar(255) NOT NULL DEFAULT 'member'
  CHECK (role in ('admin','member'));
