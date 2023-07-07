package server

import (
	"net/http"

	"github.com/TurniXXD/go-rss/utils"
)

func handleServerState(w http.ResponseWriter, r *http.Request) {
	// Empty response, just 200 status code
	utils.JSONResponse(w, 200, struct{}{})
}

func handleErr(w http.ResponseWriter, r *http.Request) {
	utils.JSONErrorResponse(w, 400, "Something went wrong")
}
