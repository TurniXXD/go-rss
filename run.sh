#! /bin/bash

source .env

if [ "$1" == "goose" ]; then
  cd sql/schema
  goose postgres $DB_URL $2
elif [ "$1" == "db-gen" ]; then
  sqlc generate
# Build & Run
elif [ "$1" == "bar" ]; then
  go build
  ./go-rss
fi