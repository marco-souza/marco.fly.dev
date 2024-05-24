package github

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gomarkdown/markdown"
)

type GitHubUser struct {
	Name   string `json:"name"`
	Bio    string `json:"bio"`
	Url    string `json:"html_url"`
	Avatar string `json:"avatar_url"`
}

func User(username, token string) GitHubUser {
	url := "/user"
	if len(username) > 0 {
		url = fmt.Sprintf("/users/%s", username)
	}

	log.Println("Loading profile", url)

	body, err := fetch(url, "GET", token)
	if err != nil {
		log.Fatalln("Error fetching profile", err)
	}

	// Parse the JSON response
	var user GitHubUser
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Fatalf("Failed to unmarshal body: %v", err)
	}

	return user
}

type GithubEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

func Emails(token string) ([]GithubEmail, error) {
	url := "/user/emails"
	log.Println("listing emails", url)

	body, err := fetch(url, "GET", token)
	if err != nil {
		return nil, err
	}

	// Parse the JSON response
	var emails []GithubEmail
	err = json.Unmarshal(body, &emails)
	if err != nil {
		log.Printf("Failed to unmarshal body: %v", err)
		return nil, err
	}

	return emails, nil
}

func Resume(url string) ([]byte, error) {
	log.Println("fetching resume", url)

	body, err := fetch(url, "GET", "")
	if err != nil {
		log.Fatalln("Error fetching resume", err)
	}

	return markdown.ToHTML(body, nil, nil), nil
}

func fetch(url, method, token string) ([]byte, error) {
	body := []byte{}
	if url[:4] != "http" {
		// if not a full url, build github api url
		url = fmt.Sprintf("%s%s", BASE_API_URL, url)
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
	// defer resp.Body.Close()
	if err != nil {
		err := fmt.Errorf("Failed to read body: %v", err)
		return body, err
	}

	return body, nil
}
