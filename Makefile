migrateup:
	migrate -path db/migration -database "postgresql://postgres:@localhost:5432/hrm_development?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:@localhost:5432/hrm_development?sslmode=disable" -verbose down

.PHONY: migrateup migratedown
