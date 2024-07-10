package github

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/marco-souza/marco.fly.dev/internal/constants"
)

type AuthCookies struct {
	*fiber.Ctx
	RefreshTokenKey string
	AccessTokenKey  string
}
type PersistedAuthTokens struct {
	RefreshToken string `cookie:"refresh_token"`
	AccessToken  string `cookie:"access_token"`
}

func (c *AuthCookies) SetAuthCookies(auth *AuthToken) {
	c.setCookie(c.AccessTokenKey, auth.AccessToken, auth.ExpiresIn)
	c.setCookie(c.RefreshTokenKey, auth.RefreshToken, auth.RefreshTokenExpiresIn)
}

func (c *AuthCookies) DeleteAuthCookies() {
	c.setCookie(c.AccessTokenKey, "", -1)
	c.setCookie(c.RefreshTokenKey, "", -1)
}

func (c *AuthCookies) GetAuthToken(name string) (*PersistedAuthTokens, error) {
	auth := &PersistedAuthTokens{}
	if err := c.CookieParser(auth); err != nil {
		return nil, err
	}
	return auth, nil
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

	logger.Info("set cookie: ", "name", name, "cookie", cookie)

	c.Cookie(cookie)
}

func AccessToken(c *fiber.Ctx) string {
	return c.Cookies(constants.ACCESS_TOKEN_KEY, "")
}

func RefreshToken(c *fiber.Ctx) string {
	return c.Cookies(constants.REFRESH_TOKEN_KEY, "")
}

func HasAccessToken(c *fiber.Ctx) bool {
	return AccessToken(c) != ""
}

func HasRefreshToken(c *fiber.Ctx) bool {
	return RefreshToken(c) != ""
}
