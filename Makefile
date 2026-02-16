postgres:
	docker run --name postgres-trixie -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:17 

createdb:
	docker exec -it postgres-trixie createdb --username=root --owner=root TeslaBank

migrate-up:
	migrate -path db/migrate -database "postgres://root:root@localhost:5432/TeslaBank?sslmode=disable" -verbose up

migrate-down:
	migrate -path db/migrate -database "postgres://root:root@localhost:5432/TeslaBank?sslmode=disable" -verbose down

dropdb:
	docker exec -it postgres-trixie dropdb TeslaBank

sqlc-gen :
	sqlc generate

test :
	go test -v -cover ./...

server :
	go run main.go

mockgen :
	cd ../..
	mv go.mod go-temp.mod
	mv go-spoof.mod go.mod
	mockgen -destination db/mock/store.go TeslaCoil196/db/sqlc Store
	mv go.mod go-spoof.mod
	mv go-temp.mod go.mod

mockgen-win :
	cd ../..
	rename go.mod go-temp.mod
	rename go-spoof.mod go.mod
	mockgen -destination db/mock/store.go TeslaCoil196/db/sqlc Store
	rename go.mod go-spoof.mod
	rename go-temp.mod go.mod

.PHONY: createdb dropdb postgres migrate-up migrate-down sqlc-gen test server mockgen-win