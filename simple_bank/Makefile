include .env

# =============================== HELPERS ===============================
# =======================================================================


.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## confirm: confirm the action
.PHONY: confirm
confirm:
	@echo -n "Are you sure? [y/N]" && read ans && [ $${ans:-N} = y ]


# =============================== TARGETS ===============================
# =======================================================================

.PHONY: pgadmin
pgadmin:
	@echo "Starting pgadmin..."
	docker run -p 50505:80 -e PGADMIN_DEFAULT_EMAIL=test@gmail.com -e PGADMIN_DEFAULT_PASSWORD=${POSTGRES_PASSWORD} --name beifong -d dpage/pgadmin4

.PHONY: postgres
postgres: confirm
	@echo "Starting postgres..."
	@echo "docker run -d --name voldermort -p ${HOST_PORT}:${DB_PORT} -e POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} postgres"

.PHONY: createdb
createdb: confirm
	@echo "Creating database..."
	docker exec -it voldermort createdb --username=${POSTGRES_USER} --owner=${POSTGRES_USER} $(DB_NAME) 

.PHONY: dropdb
dropdb: confirm
	@echo "Dropping database..."
	docker exec -it voldermort dropdb ${DB_NAME}

.PHONY: format
format:
	@echo "Formatting code..."
	find . -name "*.go" | xargs gofmt -w

.PHONY: migrateup
migrateup: confirm
	@echo "Migrating up..."
	migrate -path db/migration -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose up

.PHONY: migratedown
migratedown: confirm
	@echo "Migrating down..."
	migrate -path db/migration -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose down

.PHONY: sqlc
sqlc:
	@echo "Generating SQL code..."
	sqlc generate

.PHONY: test
test:
	go test -v -cover -count=1 -timeout 30s ./...	

