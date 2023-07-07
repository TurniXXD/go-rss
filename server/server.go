package server

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/TurniXXD/go-rss/internal/database"
	"github.com/TurniXXD/go-rss/scraper"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

type apiConfig struct {
	DB *database.Queries
}

var apiCfg apiConfig

func Init() {
	db := initDb()

	portString := os.Getenv("PORT")
	if portString == "" {
		portString = "8080"
	}

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	srv := &http.Server{
		Handler: r,
		Addr:    ":" + portString,
	}

	go scraper.Start(
		db,
		10,
		time.Minute,
	)

	v1r := chi.NewRouter()
	v1r.Get("/health", handleServerState)
	v1r.Get("/err", handleErr)

	v1r.Post("/users", apiCfg.handleCreateUser)
	v1r.Get("/users", apiCfg.middlewareAuth(apiCfg.handleGetUser))

	v1r.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handleCreateFeed))
	v1r.Get("/feeds", apiCfg.handleGetFeeds)

	v1r.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handleCreateFeedFollow))
	v1r.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handleGetFeedFollows))
	v1r.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handleDeleteFeedFollow))

	v1r.Get("/posts", apiCfg.middlewareAuth(apiCfg.handleGetPostsForUser))

	r.Mount("/v1", v1r)

	time.AfterFunc(time.Second*3, func() {
		log.Printf("Server running on port %s", portString)
	})

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
