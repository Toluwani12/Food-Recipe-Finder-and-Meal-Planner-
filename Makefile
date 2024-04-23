install-migrate-mac:
	curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.darwin-amd64.tar.gz | tar xvz
	mv migrate.darwin-amd64 /usr/local/bin/migrate

install-migrate-windows:
	curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.windows-amd64.exe.tar.gz | tar xvz
	mv migrate.windows-amd64.exe /usr/local/bin/migrate

migrate-up-all:
	migrate -path db/migrations -database "postgresql://localhost:5432/recipe?sslmode=disable" up

migrate-up:
	migrate -path db/migrations -database "postgresql://localhost:5432/recipe?sslmode=disable" up 1

migrate-down-all:
	migrate -path db/migrations -database "postgresql://localhost:5432/recipe?sslmode=disable" down

migrate-down:
	migrate -path db/migrations -database "postgresql://localhost:5432/recipe?sslmode=disable" down 1

migrate-force:
	@if [ -z "$(version)" ]; then echo "version is not set. Set it like this: make migrate-force version=4"; exit 1; fi
	@migrate -path db/migrations -database "postgresql://localhost:5432/recipe?sslmode=disable" force $(version)

migration:
	@if [ -z "$(name)" ]; then echo "name is not set. Set it like this: make migration name=create_users"; exit 1; fi
	@migrate create -ext sql -dir db/migrations $(name)

run :
	go run main.go

build :
	go build -o bin/main main.go

test :
	go test -v ./...