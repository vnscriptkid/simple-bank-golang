DB_URL=postgresql://postgres:123456@localhost:5433/bank?sslmode=disable

network:
	docker network create bank-network

pg:
	docker run --name bank-pg --network bank-network -p 5433:5432 -e POSTGRES_USER=postgres -e POSTGRES_DB=bank -e POSTGRES_PASSWORD=123456 -d postgres:14.5

createdb:
	docker exec -it bank-pg createdb --username=postgres --owner=postgres bank

dropdb:
	docker exec -it bank-pg dropdb --username=postgres bank

migratecreate:
	migrate create -ext sql -dir db/migration -seq init_schema

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...