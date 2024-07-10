package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"
	"strings"

	"github.com/marco-souza/marco.fly.dev/internal/config"
)

var (
	logger           = slog.With("service", "github")
	authUrl          = "https://github.com/login/oauth"
	contentType      = "application/json"
	cfg              = config.Load()
	githuCredentials = GithubCredential{
		ClientId:     cfg.Github.ClientId,
		ClientSecret: cfg.Github.ClientSecret,
	}
)

type Auth struct {
	AllowedEmails map[string]bool
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
	emails, err := Emails(token)
	if err != nil {
		logger.Error("failed to fetch emails", "err", err)
		return false
	}

	logger.Info("checking user emails", "emails", emails, "allowed", a.AllowedEmails)

	for _, email := range emails {
		isValid, ok := a.AllowedEmails[email.Email]
		if !ok {
			logger.Warn("email is not allowed, checking next")
			continue
		}

		if isValid {
			return true
		}
	}

	logger.Info("user is not allowed")
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
		logger.Warn("failed to parse url", "url", authUrl)
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

	logger.Info("parsed url", "url", redirectUrl.String())

	return redirectUrl.String()
}

func fetchAuthToken(params interface{}) (*AuthToken, error) {
	logger.Info("fetching auth token", "params", params)
	body, err := json.Marshal(params)
	if err != nil {
		logger.Error("failed to parse params", "params", params, "err", err)
		return nil, err
	}

	logger.Info("fetch auth token with code with credentials")

	client := &http.Client{}
	accessTokenUrl := authUrl + "/access_token"
	req, err := http.NewRequest("POST", accessTokenUrl, bytes.NewBuffer(body))
	if err != nil {
		logger.Error("failed to fetch access token", "err", err)
		return nil, err
	}

	req.Header.Set("Accept", contentType)
	req.Header.Set("Content-Type", contentType)

	res, err := client.Do(req)
	if err != nil {
		logger.Error("failed to fetch access token", "err", err)
		return nil, err
	}

	logger.Info("read body buffer")
	accessTokenRaw, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		logger.Error("failed to parse body", "err", err)
		return nil, err
	}

	logger.Info("parse json response", "accessTokenRaw", accessTokenRaw)
	var accessToken AuthToken
	err = json.Unmarshal(accessTokenRaw, &accessToken)
	if err != nil {
		logger.Error("failed to unmarshal body", "err", err)
		return nil, err
	}

	logger.Info("access token is here", "accessToken", accessToken)

	return &accessToken, nil
}
