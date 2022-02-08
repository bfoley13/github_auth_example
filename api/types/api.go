package types

type APIError struct {
	Error string `json:"error"`
}

type Response struct {
	Data interface{} `json:"data"`
}

type OAuthURLResponse struct {
	URL string `json:"url"`
}

type OAuthTokenResponse struct {
	Token string `json:"token"`
}

type GithubAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}
