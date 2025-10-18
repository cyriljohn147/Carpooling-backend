postgres:
	docker run --name some-postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres

startdb:
	docker start some-postgres

stopdb:
	docker stop some-postgres

createdb:
	docker exec -it some-postgres createdb --username=root --owner=root carpooling

createdb-ci:
	docker exec some-postgres createdb --username=root --owner=root carpooling

dropdb:
	docker exec -it some-postgres dropdb carpooling

dropdb-ci:
	docker exec some-postgres dropdb carpooling

migrateup:
	migrate -path db/migration/ -database "postgresql://root:secret@localhost:5432/carpooling?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration/ -database "postgresql://root:secret@localhost:5432/carpooling?sslmode=disable" -verbose down

runpsql:
	docker exec -it some-postgres psql -U root

sqlc:
	sqlc generate

createtestdb:
	docker exec -it some-postgres createdb --username=root --owner=root carpooling_test

createtestdb-ci:
	docker exec some-postgres createdb --username=root --owner=root carpooling_test

droptestdb:
	docker exec -it some-postgres dropdb carpooling_test

droptestdb-ci:
	docker exec some-postgres dropdb carpooling_test

migrateup-test:
	migrate -path db/migration/ -database "postgresql://root:secret@localhost:5432/carpooling_test?sslmode=disable" -verbose up

migratedown-test:
	migrate -path db/migration/ -database "postgresql://root:secret@localhost:5432/carpooling_test?sslmode=disable" -verbose down

test:
	go test -v -coverpkg=./db/sqlc -coverprofile=coverage.out ./db/sqlc/test

test-coverage:
	go test -v -coverpkg=./db/sqlc -coverprofile=coverage.out ./db/sqlc/test
	go tool cover -func=coverage.out

test-coverage-html:
	go test -v -coverpkg=./db/sqlc -coverprofile=coverage.out ./db/sqlc/test
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-users:
	go test -v ./db/sqlc/test -run TestUser

test-trips:
	go test -v ./db/sqlc/test -run TestTrip

test-participants:
	go test -v ./db/sqlc/test -run TestParticipant

test-balances:
	go test -v ./db/sqlc/test -run TestBalance

setup-test: createtestdb migrateup-test
	@echo "Test database setup complete"

teardown-test: droptestdb
	@echo "Test database torn down"

# Server commands
build:
	go build -o bin/carpooling-api main.go

run:
	go run main.go

dev:
	go run main.go

clean:
	rm -rf bin/
	go clean

# Install dependencies
deps:
	go mod tidy
	go mod download

# Code quality
fmt:
	go fmt ./...

vet:
	go vet ./...

check: fmt vet test
	@echo "All checks passed"

.PHONY: postgres startdb stopdb createdb createdb-ci dropdb dropdb-ci migrateup migratedown sqlc runpsql createtestdb createtestdb-ci droptestdb droptestdb-ci migrateup-test migratedown-test test test-coverage test-coverage-html test-users test-trips test-participants test-balances setup-test teardown-test build run dev clean deps fmt vet check
