CREATE TABLE IF NOT EXISTS leave_types (
  id bigserial PRIMARY KEY,
  name varchar(255) NOT NULL,
  usage int NOT NULL DEFAULT 0,
  max int NOT NULL,
  user_id int NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL,
  CONSTRAINT fk_leave_types_users FOREIGN KEY (user_id) REFERENCES users (id)
);
