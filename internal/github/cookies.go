package github

import (
	"strings"

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

func (c *AuthCookies) setCookie(name, token string, expiry int) {
	cookie := fiber.Cookie{
		HTTPOnly: true,
		Path:     "/",
		Name:     name,
		Value:    token,
		MaxAge:   expiry,
		Domain:   c.Hostname(),
		Secure:   strings.HasPrefix(c.Protocol(), "https"),
	}

	c.Cookie(&cookie)
}

func (c *AuthCookies) DeleteAuthCookies() {
	c.ClearCookie(c.RefreshTokenKey, c.AccessTokenKey)
}
