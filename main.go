package main

import (
	"github_auth_app/api/routes"
)

func main() {
	gitHubOAuthService := routes.NewGitHubOAuthService()
	gitHubOAuthService.StartApiService(routes.NewRouter(gitHubOAuthService))
}
