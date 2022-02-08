package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(apiServer *GitHubOAuthService) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/login/callback", corsHandler(apiServer.GetAuthToken, "GET")).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/login", corsHandler(apiServer.GetGitHubOAuthURL, "GET")).Methods("GET", "OPTIONS")
	return router
}

func corsHandler(h func(http.ResponseWriter, *http.Request), validMethods string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", validMethods)
		} else {
			h(w, r)
		}
	}
}
