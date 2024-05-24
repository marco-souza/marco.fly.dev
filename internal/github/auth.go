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
	"strings"

	"github.com/marco-souza/marco.fly.dev/internal/config"
)

var (
	authUrl          = "https://github.com/login/oauth"
	contentType      = "application/json"
	cfg              = config.Load()
	githuCredentials = GithubCredential{
		ClientId:     cfg.Github.ClientId,
		ClientSecret: cfg.Github.ClientSecret,
	}
)

type Auth struct {
	AllowedEmails []string
}

type GithubCredential struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type RefreshAuthParam struct {
	*GithubCredential
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
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

func (a *Auth) IsUserAllowed(token string) bool {
	emails := Emails(token)
	if emails == nil {
		log.Println("failed to fetch emails")
		return false
	}

	log.Printf("emails: %v", emails)

	for _, email := range *emails {
		for _, allowedEmail := range a.AllowedEmails {
			if email.Email == allowedEmail {
				log.Printf("user is allowed - email: %v, allowed: %v", email, a.AllowedEmails)
				return true
			}
		}
	}

	return false
}

func (a *Auth) FetchAuthToken(code string) (*AuthToken, error) {
	return fetchAuthToken(
		FetchAuthParam{
			Code:             code,
			GithubCredential: &githuCredentials,
		},
	)
}

func (a *Auth) RefreshAuthToken(refreshToken string) (*AuthToken, error) {
	return fetchAuthToken(
		RefreshAuthParam{
			GrantType:        "refresh_token",
			RefreshToken:     refreshToken,
			GithubCredential: &githuCredentials,
		},
	)
}

func (a *Auth) RedirectLink(origin string) string {
	redirectUrl, err := url.Parse(authUrl + "/authorize")
	if err != nil {
		log.Fatalf("failed to parse url: %s", authUrl)
	}

	query := url.Values{}
	state := fmt.Sprintf("%d", rand.Intn(1000000000))
	redirectUri := origin + cfg.Github.CallbackUrl

	query.Add("response_type", "code")
	query.Add("state", state)
	query.Add("scope", cfg.Github.Scope)
	query.Add("client_id", cfg.Github.ClientId)
	query.Add("redirect_uri", redirectUri)

	redirectUrl.RawQuery = strings.ReplaceAll(
		query.Encode(), "+", "%20",
	)

	log.Printf("parsed url: %s", redirectUrl.String())

	return redirectUrl.String()
}

func fetchAuthToken(params interface{}) (*AuthToken, error) {
	log.Print("parse params", params)
	body, err := json.Marshal(params)
	if err != nil {
		log.Fatal("failed to parse params", params, err)
		return nil, err
	}

	log.Print("fetch auth token with code + credentials")

	client := &http.Client{}
	accessTokenUrl := authUrl + "/access_token"
	req, err := http.NewRequest("POST", accessTokenUrl, bytes.NewBuffer(body))
	if err != nil {
		log.Fatal("failed to fetch access token", err)
		return nil, err
	}

	req.Header.Set("Accept", contentType)
	req.Header.Set("Content-Type", contentType)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal("failed to fetch access token", err)
		return nil, err
	}

	log.Print("read body buffer")
	accessTokenRaw, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Fatalf("failed to parse body: %v", err)
		return nil, err
	}

	log.Print("parse json response ", string(accessTokenRaw))
	var accessToken AuthToken
	err = json.Unmarshal(accessTokenRaw, &accessToken)
	if err != nil {
		log.Fatalf("failed to unmarshal body: %v", err)
		return nil, err
	}

	log.Print("access token is here: ", accessToken)

	return &accessToken, nil
}
