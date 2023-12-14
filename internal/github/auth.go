package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/marco-souza/marco.fly.dev/internal/config"
)

var (
	authUrl          = BASE_API_URL + "/oauth"
	contentType      = "application/json"
	cfg              = config.Load()
	githuCredentials = GithubCredential{
		ClientId:     cfg.Github.ClientId,
		ClientSecret: cfg.Github.ClientSecret,
	}
)

type Auth struct{}

type GithubCredential struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type FetchAuthParam struct {
	*GithubCredential
	Code string `json:"code"`
}

type AuthToken struct {
	Scope                 string `json:"scope"`
	TokenType             string `json:"token_type"`
	AccessToken           string `json:"access_token"`
	RefreshToken          string `json:"refresh_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
}

func (a *Auth) FetchAuthToken(code string) (*AuthToken, error) {
	params := FetchAuthParam{
		Code:             code,
		GithubCredential: &githuCredentials,
	}
	}

	// parse params
	body, err := json.Marshal(params)
	if err != nil {
		log.Fatal("failed to parse params", params, err)
		return nil, err
	}

	// fetch auth token with code + credentials
	accessTokenUrl := authUrl + "/access_token"
	res, err := http.Post(accessTokenUrl, contentType, bytes.NewBuffer(body))
	if err != nil {
		log.Fatal("failed to fetch access token", err)
		return nil, err
	}

	// read body buffer
	accessTokenRaw, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Fatalf("Failed to parse body: %v", err)
		return nil, err
	}

	// Parse the JSON response
	var accessToken AuthToken
	err = json.Unmarshal(accessTokenRaw, &accessToken)
	if err != nil {
		log.Fatalf("Failed to unmarshal body: %v", err)
		return nil, err
	}

	log.Print("access token is here: ", accessToken)

	return &accessToken, nil
}

func (a *Auth) RedirectLink(origin string) string {
	redirectUrl, err := url.Parse(authUrl + "/authorize")
	if err != nil {
		log.Fatalf("Failed to parse url: %s", authUrl)
	}

	query := url.Values{}
	state := fmt.Sprintf("%d", rand.Intn(1000000000))
	redirectUri := origin + cfg.Github.CallbackUrl

	query.Add("response_type", "code")
	query.Add("state", state)
	query.Add("scope", cfg.Github.Scope)
	query.Add("client_id", cfg.Github.ClientId)
	query.Add("client_secret", cfg.Github.ClientSecret)
	query.Add("redirect_uri", redirectUri)

	redirectUrl.RawQuery = query.Encode()

	log.Printf("Failed to parse url: %s", redirectUrl.String())
	return redirectUrl.String()
}
