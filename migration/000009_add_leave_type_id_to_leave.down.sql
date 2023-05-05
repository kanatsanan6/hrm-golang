ALTER TABLE "leaves"
  DROP CONSTRAINT fk_leaves_leave_types;

ALTER TABLE "leaves"
  DROP leave_type_id;
