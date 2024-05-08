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
	// set up the GitHub API endpoint
	url := fmt.Sprintf("%s/users/%s", BASE_API_URL, username)
	if len(username) == 0 {
		log.Println("Loading logged user profile", username, token)
		url = fmt.Sprintf("%s/user", BASE_API_URL)
	}

	log.Println("Loading profile", url)

	// make a GET request to the URL
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	if len(token) > 0 {
		req.Header.Add("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to retrieve profile data: %d", resp.StatusCode)
	}

	// read body
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatalf("Failed to parse body: %v", err)
	}

	// Parse the JSON response
	var user GitHubUser
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Fatalf("Failed to unmarshal body: %v", err)
	}

	return user
}

func Resume(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching resume", err)
		return nil, err
	}

	// parse resume body into a html template
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading resume body", err)
		return nil, err
	}

	return markdown.ToHTML(body, nil, nil), nil
}
