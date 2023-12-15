package github

import (
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthCookies struct {
	*fiber.Ctx
	RefreshTokenKey string
	AccessTokenKey  string
}

func (c *AuthCookies) SetAuthCookies(auth *AuthToken) {
	c.setCookie(c.AccessTokenKey, auth.AccessToken, auth.ExpiresIn)
	c.setCookie(c.RefreshTokenKey, auth.RefreshToken, auth.RefreshTokenExpiresIn)
}

func (c *AuthCookies) DeleteAuthCookies() {
	c.setCookie(c.AccessTokenKey, "", -1)
	c.setCookie(c.RefreshTokenKey, "", -1)
}

func (c *AuthCookies) setCookie(name, token string, expires int) {
	cookie := &fiber.Cookie{
		Name:        name,
		Value:       token,
		Path:        "/",
		HTTPOnly:    true,
		SameSite:    "Lax",
		SessionOnly: false,
		Expires:     time.Now().Add(time.Second * time.Duration(expires)),
		Domain:      strings.Join(c.GetReqHeaders()["host"], ""),
		Secure:      strings.HasPrefix(c.Protocol(), "https"),
	}

	log.Print("set cookie: ", name, cookie)

	c.Cookie(cookie)
}
