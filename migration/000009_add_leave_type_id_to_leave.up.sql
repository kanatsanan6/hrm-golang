ALTER TABLE "leaves"
  ADD leave_type_id int NOT NULL,
  DROP leave_type;

ALTER TABLE "leaves" ADD CONSTRAINT fk_leaves_leave_types FOREIGN KEY (leave_type_id) REFERENCES leave_types (id);
