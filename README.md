# go-rss

- [Full course](https://youtu.be/un6ZyFkqFKo?t=24196)
- Simple REST API for creating new user and assigning API key to him
- User can then get his user info and create new feeds

## Commands

- build & run => `go build && ./go-rss`
- run migrations => `cd sql/schema && goose postgres postgres://user:password@localhost:5432/go-rss-db up`
- rollback single migration => `cd sql/schema && goose postgres postgres://user:password@localhost:5432/go-rss-db down`
- gen sqlc => `sqlc generate`

### Run script

- `./run.sh ARG_1 ARG_2`
- supported args: `goose up`, `goose down`, `bar`, `db-gen`

## Tech

- chi router
- sqlc => `go install github.com/kyleconroy/sqlc/cmd/sqlc@latest`
- goose => `go install github.com/pressly/goose/v3/cmd/goose@latest`
