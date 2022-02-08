package main

import (
	"github_auth_app/api/routes"
)

func main() {
	topShelfService := routes.NewGitHubOAuthService()
	topShelfService.StartApiService(routes.NewRouter(topShelfService))
}
