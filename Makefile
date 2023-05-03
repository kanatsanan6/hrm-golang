migrateup:
	migrate -path migration -database "postgresql://postgres:secret@localhost:5432/hrm_development?sslmode=disable" -verbose up

migratedown:
	migrate -path migration -database "postgresql://postgres:secret@localhost:5432/hrm_development?sslmode=disable" -verbose down

mockgen_queries:
	mockgen --destination queries/mock/queries.go github.com/kanatsanan6/hrm/queries Queries

mockgen_service:
	mockgen --destination service/mock/service.go github.com/kanatsanan6/hrm/service Service

.PHONY: migrateup migratedown mockgen_queries mockgen_service
