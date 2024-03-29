# Pickup API

## Getting Started

- go >= 1.21.x [brew](https://formulae.brew.sh/formula/go) or golang [official](https://go.dev/doc/install)
- swag `go install github.com/swaggo/swag/cmd/swag@latest`
  - to regenerate openapi docs `swag init`
- gow `go install github.com/mitranim/gow@latest`
- install [sqlx-cli](https://github.com/launchbadge/sqlx/tree/main/sqlx-cli)
- start local db `docker compose up -d`
- [swagger](https://petstore.swagger.io/?url=http://localhost:3000/openapi.json)

## Migrations

- create the database `sqlx database create`
- run migrations `sqlx migrate run`
- prepare for offline `cargo sqlx prepare`

- drop the database `sqlx database drop`
- revert migrations `sqlx migrate revert`
