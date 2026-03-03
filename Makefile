postgres:
	docker run --name postgres-trixie --network tesla-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:17 

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

build-network:
	docker network create tesla-network

buildbank:
	docker build -t teslabank:latest .

runbank:
	Docker run --name teslabank --network tesla-network -p 8080:8080 -e DB_SOURCE='postgres://root:root@postgres-trixie:5432/TeslaBank?sslmode=disable' teslabank:latest

proto:
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto

evens:
	evans --host localhost --port 9090 -r repl --reflection
	show package
	package pb
	show service

.PHONY: createdb dropdb postgres migrate-up migrate-down sqlc-gen test server mockgen-win proto