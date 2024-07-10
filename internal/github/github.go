package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gomarkdown/markdown"
	"github.com/marco-souza/marco.fly.dev/internal/cache"
)

type GitHubUser struct {
	Name   string `json:"name"`
	Bio    string `json:"bio"`
	Url    string `json:"html_url"`
	Avatar string `json:"avatar_url"`
}

func User(username, token string) (*GitHubUser, error) {
	url := "/user"
	if len(username) > 0 {
		url = fmt.Sprintf("/users/%s", username)
	}

	logger.Info("loading profile")

	body, err := fetch(url, "GET", token)
	if err != nil {
		logger.Error("error fetching profile", "username", username, "err", err)
		return nil, err
	}

	// Parse the JSON response
	var user GitHubUser
	err = json.Unmarshal(body, &user)
	if err != nil {
		logger.Error("Failed to unmarshal body", "err", err)
		return nil, err
	}

	return &user, nil
}

type GithubEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

func Emails(token string) ([]GithubEmail, error) {
	url := "/user/emails"
	logger.Info("listing emails", "url", url, "token", token)

	body, err := fetch(url, "GET", token)
	if err != nil {
		return nil, err
	}

	// Parse the JSON response
	var emails []GithubEmail
	err = json.Unmarshal(body, &emails)
	if err != nil {
		logger.Error("Failed to unmarshal body", "err", err)
		return nil, err
	}

	return emails, nil
}

func Resume(url string) ([]byte, error) {
	logger.Info("fetching resume", "url", url)

	body, err := fetch(url, "GET", "")
	if err != nil {
		logger.Error("Error fetching resume", "err", err)
		return nil, err
	}

	return markdown.ToHTML(body, nil, nil), nil
}

func fetch(url, method, token string) ([]byte, error) {
	body := []byte{}
	if url[:4] != "http" {
		// if not a full url, build github api url
		url = fmt.Sprintf("%s%s", BASE_API_URL, url)
	}

	cacheKey := fmt.Sprintf("%s %s", method, url)
	if cached, err := cache.Get(cacheKey); err != nil {
		logger.Info("miss", "key", cacheKey)
	} else {
		logger.Info("hit", "key", cacheKey)
		return cached, nil
	}

	// make a GET request to the URL
	client := http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return body, err
	}

	if len(token) > 0 {
		req.Header.Add("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return body, err
	}

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("Failed to retrieve data: %d", resp.StatusCode)
		return body, err
	}

	// read body
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		err := fmt.Errorf("Failed to read body: %v", err)
		return body, err
	}

	// cache content by 15 minutes
	cache.Set(cacheKey, body, cache.WithTTL(15*60))

	return body, nil
}
