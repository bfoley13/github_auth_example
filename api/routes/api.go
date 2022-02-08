package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github_auth_app/api/types"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type GitHubOAuthService struct{}

func NewGitHubOAuthService() *GitHubOAuthService {
	return &GitHubOAuthService{}
}

func (api *GitHubOAuthService) StartApiService(router *mux.Router) {
	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	srv.Shutdown(ctx)

	log.Println("shutting down")
	os.Exit(0)
}

func (api *GitHubOAuthService) WriteHTTPErrorResponse(w http.ResponseWriter, code int, errResp error) {
	log.Println("[WriteHTTPErrorResponse] Error: ", errResp.Error())
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(&types.APIError{Error: errResp.Error()}); err != nil {
		log.Println("[WriteHTTPErrorResponse] failed to write http response: ", err.Error())
	}
}

func getGitHubAppClientId() string {
	appClientId, ok := os.LookupEnv("APP_CLIENT_ID")
	if !ok {
		log.Println("[getGitHubAppClientId] failed to find app client id in environment variable")
		return ""
	}

	return appClientId
}

func getGitHubAppSecret() string {
	appSecret, ok := os.LookupEnv("APP_SECRET")
	if !ok {
		log.Println("[getGitHubAppSecret] failed to find app secret in environment variable")
		return ""
	}

	return appSecret
}

func getCallBackURI() string {
	appRedirectURI, ok := os.LookupEnv("APP_REDIRECT_URI")
	if !ok {
		log.Println("[getCallBackURI] failed to find redirect uri in environment variable")
		return ""
	}

	return appRedirectURI
}

func (api *GitHubOAuthService) GetGitHubOAuthURL(w http.ResponseWriter, r *http.Request) {
	log.Println("[GetGitHubOAuthURL] starting service call")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	responseURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&state=%s&scope=%s",
		getGitHubAppClientId(),
		getCallBackURI(),
		uuid.NewV4(),
		"repo,user:email,workflow",
	)

	log.Println("[GetGitHubOAuthURL] oauthURI: ", responseURL)
	if err := json.NewEncoder(w).Encode(types.OAuthURLResponse{URL: responseURL}); err != nil {
		api.WriteHTTPErrorResponse(w, 500, err)
		return
	}

	log.Println("[GetGitHubOAuthURL] finished")
	w.WriteHeader(http.StatusOK)
}

func (api *GitHubOAuthService) GetAuthToken(w http.ResponseWriter, r *http.Request) {
	log.Println("[GetAuthToken] starting service call")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	code := r.URL.Query().Get("code")

	githubAuthToken, err := getGithubAuthToken(code)
	if err != nil {
		log.Println("[GetAuthToken] failed to get auth token: ", err.Error())
		api.WriteHTTPErrorResponse(w, 500, fmt.Errorf("failed to get quth token"))
	}

	if err := json.NewEncoder(w).Encode(types.OAuthTokenResponse{Token: githubAuthToken}); err != nil {
		api.WriteHTTPErrorResponse(w, 500, err)
		return
	}

	log.Println("[GetAuthToken] finished")
	w.WriteHeader(http.StatusOK)
}

func getGithubAuthToken(code string) (string, error) {
	clientID := getGitHubAppClientId()
	clientSecret := getGitHubAppSecret()

	log.Println("[getGithubAuthToken] building new request")
	requestBodyMap := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"code":          code,
		"redirect_uri":  getCallBackURI(),
	}
	requestJSON, _ := json.Marshal(requestBodyMap)

	log.Println("[getGithubAuthToken] sending request")
	req, reqerr := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBuffer(requestJSON),
	)
	if reqerr != nil {
		log.Println("[getGithubAuthToken] failed to mkae request: ", reqerr.Error())
		return "", reqerr
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Println("[getGithubAuthToken] request failed: ", resperr.Error())
		return "", resperr
	}

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[getGithubAuthToken] failed to parse response body: ", err.Error())
		return "", err
	}

	gitHubAuthResponse := types.GithubAccessTokenResponse{}
	json.Unmarshal(respbody, &gitHubAuthResponse)

	return gitHubAuthResponse.AccessToken, nil
}
