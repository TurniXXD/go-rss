package main

import (
	"github.com/TurniXXD/go-rss/server"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	server.Init()
}
