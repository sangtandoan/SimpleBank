postgres:
	docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

createdb: 
	docker exec -it postgres16 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres16 dropdb simple_bank

migrateup:
	migrate -path db/migrate -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 

# migrate up 1 version
migrateup1:
	migrate -path db/migrate -database "postgresql://root:yBjq2bhYW2YJQ7Sh45p1@localhost:5432/simple_bank?sslmode=disable" -verbose up 1
	
migratedown:
	migrate -path db/migrate -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

# migrate down 1 version
migratedown1:
	migrate -path db/migrate -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

test:
	go test -v ./...

testcover:
	go test ./... -coverprofile=c.out -covermode=count

coverhtml:
	go tool cover -html=c.out



server:
	go run ./cmd/main.go

mock: 
	~/go/bin/mockgen -package mockdb -destination internal/mock/store.go github.com/FrostJ143/simplebank/internal/query Store


.PHONY: createdb dropdb postgres migrateup migratedown test server mock testcover coverhtml migrateup1 migratedown1
