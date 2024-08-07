MIGRATE_VERSION=v4.14.1
MIGRATE_URL=https://github.com/golang-migrate/migrate/releases/download/$(MIGRATE_VERSION)
MIGRATE_MAC=migrate.darwin-amd64
MIGRATE_WINDOWS=migrate.windows-amd64.exe
MIGRATE_BIN=/usr/local/bin/migrate

DB_URL=postgres://recipe_gc2u_user:qCLSFmKiOGfROBK3OzfjwYmsqGMRF5yQ@dpg-cqd760bv2p9s73e95o6g-a.frankfurt-postgres.render.com:5432/recipe_gc2u?sslmode=require
#DB_URL=postgres://recipe_q9i7_user:HcRB9IqNE49we2aN93XGp7z5RPdZY0rq@dpg-cpcf6t3tg9os738ccrjg-a.frankfurt-postgres.render.com:5432/recipe_q9i7?sslmode=require
#DB_URL=postgresql://localhost:5432/recipe?sslmode=disable
MIGRATE_CMD=migrate -path db/migrations -database "$(DB_URL)"

install-migrate-mac:
	curl -L $(MIGRATE_URL)/$(MIGRATE_MAC).tar.gz | tar xvz
	mv $(MIGRATE_MAC) $(MIGRATE_BIN)

install-migrate-windows:
	curl -L $(MIGRATE_URL)/$(MIGRATE_WINDOWS).tar.gz | tar xvz
	mv $(MIGRATE_WINDOWS) $(MIGRATE_BIN)

migrate-up-all:
	$(MIGRATE_CMD) up

migrate-up:
	$(MIGRATE_CMD) up 1

migrate-down-all:
	$(MIGRATE_CMD) down

migrate-down:
	$(MIGRATE_CMD) down 1

migrate-force:
	@if [ -z "$(version)" ]; then echo "version is not set. Set it like this: make migrate-force version=4"; exit 1; fi
	@migrate -path db/migrations -database "postgresql://localhost:5432/recipe?sslmode=disable" force $(version)

migration:
	@if [ -z "$(name)" ]; then echo "name is not set. Set it like this: make migration name=create_users"; exit 1; fi
	@migrate create -ext sql -dir db/migrations $(name)

run:
	go run main.go

build:
	go build -o bin/main main.go

test:
	go test -v ./...
