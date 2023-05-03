ALTER TABLE "leaves"
  ADD user_id int NOT NULL;

ALTER TABLE "leaves" ADD CONSTRAINT fk_leaves_users FOREIGN KEY (user_id) REFERENCES users (id);
