package github

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"

	"github.com/marco-souza/marco.fly.dev/internal/config"
)

const authUrl = BASE_API_URL + "/oauth"

var cfg = config.Load()

type Auth struct{}

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
