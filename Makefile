install-migrate-mac:
	curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.darwin-amd64.tar.gz | tar xvz
	mv migrate.darwin-amd64 /usr/local/bin/migrate

install-migrate-windows:
	curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.windows-amd64.exe.tar.gz | tar xvz
	mv migrate.windows-amd64.exe /usr/local/bin/migrate

migrate-up:
	migrate -path db/migrations -database "postgresql://localhost:5432/recipe?sslmode=disable" up

migrate-down:
	migrate -path db/migrations -database "postgresql://localhost:5432/recipe?sslmode=disable" down

migrate-force:
	@if [ -z "$(version)" ]; then echo "version is not set. Set it like this: make migrate-force version=4"; exit 1; fi
	@migrate -path db/migrations -database "postgresql://localhost:5432/recipe?sslmode=disable" force $(version)