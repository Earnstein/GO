include .env

# ==================================================================================== #
# HELPERS
# ==================================================================================== #


## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'


confirm:
	@echo -n "Are you sure? [y/N]" && read ans && [ $${ans:-N} = y ]


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #


## run: run the cmd/api application
run:
	go run ./cmd/api -port=8000


## psql: connect to the database using pgcli
psql:
	pgcli ${DB_DSN}


## migrations name=$1: create a new database migration
migration:
	@echo 'Creating migration files for ${name}'
	migrate create -seq -ext=.sql -dir=./migrations ${name} -verbose


## up: apply all up database migrations
up: confirm
	@echo "Running up migrations..."
	migrate -path ./migrations -database ${DB_DSN} -verbose up


## down: apply all down database migrations
down: confirm
	@echo "Running up migrations..."
	migrate -path ./migrations -database ${DB_DSN} -verbose down


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #


## audit: tidy dependencies and format, vet and test all code
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...



## vendor: tidy and vendor dependencies
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor


# ==================================================================================== #
# BUILD
# ==================================================================================== #


## build/api: build the cmd/api application
build:
	@echo 'Building cmd/api...'
	go build -ldflags='-s' -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/api ./cmd/api


.PHONY: run psql migration up down audit vendor build