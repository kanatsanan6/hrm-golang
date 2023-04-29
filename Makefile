migrateup:
	migrate -path migration -database "postgresql://postgres:@localhost:5432/hrm_development?sslmode=disable" -verbose up

migratedown:
	migrate -path migration -database "postgresql://postgres:@localhost:5432/hrm_development?sslmode=disable" -verbose down

mockgen:
	mockgen --destination queries/mock/queries.go github.com/kanatsanan6/hrm/queries Queries

.PHONY: migrateup migratedown mockgen
