migrateup:
	migrate -path migration -database "postgresql://postgres:secret@localhost:5432/hrm_development?sslmode=disable" -verbose up

migratedown:
	migrate -path migration -database "postgresql://postgres:secret@localhost:5432/hrm_development?sslmode=disable" -verbose down

genqueries:
	mockgen --destination queries/mock/queries.go github.com/kanatsanan6/hrm/queries Queries

genpolicy:
	mockgen --destination service/mock/policy.go github.com/kanatsanan6/hrm/service PolicyInterface

.PHONY: migrateup migratedown genqueries genpolicy
