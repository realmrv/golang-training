# Golang training project

## Requirements

- db: `MySQL/MariaDB`
- golang version: `16`
- migrate tool: [golang-migrate/migrate](https://github.com/golang-migrate/migrate)
- optional to start db server quickly: [Lando](https://github.com/lando/lando)

## Get started

The commands given in the example must be executed from module root directory

1. Setup `.env` from `.env.example`
2. Run migrations. Example:
   ```sh
    migrate -database 'mysql://mariadb:mariadb@tcp(127.0.0.1:3306)/database' -path db/migrations up
   ```
3. Install all project dependencies: `go install`
3. Compile and start: `go run main.go`
4. Enjoy. Default address: `http://localhost:3000/`
