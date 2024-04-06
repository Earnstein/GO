
POSTGRES_USER = root
POSTGRES_PASSWORD = secret
DB_NAME = simple_bank
DB_PORT = 5432
HOST_PORT = 5432

# Targets

postgres:
	docker run -d --name voldermort -p $(HOST_PORT):$(DB_PORT) -e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) postgres

createdb:
	docker exec -it voldermort createdb --username=$(POSTGRES_USER) --owner=$(POSTGRES_USER) $(DB_NAME) 

dropdb:
	docker exec -it voldermort dropdb $(DB_NAME)

migrateup:
	migrate -path db/migration -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover -count=1 -timeout 30s ./...	
.PHONY: postgres createdb dropdb migrateup migratedown
