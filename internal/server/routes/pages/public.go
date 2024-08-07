package pages

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/marco-souza/marco.fly.dev/internal/config"
	"github.com/marco-souza/marco.fly.dev/internal/github"
)

var logger = slog.With("service", "pages")

type rootProps struct {
	config.PageParams
	PrimaryBtn   string
	SecondaryBtn string
	Profile      github.GitHubUser
	Description  template.HTML
}

func notFoundHandler(c *fiber.Ctx) error {
	logger.Info("page not found")
	return c.
		Status(http.StatusNotFound).
		Render("404", config.DefaultPageParams)
}

func rootHandler(c *fiber.Ctx) error {
	user, _ := github.User("marco-souza", "")
	props := rootProps{
		PageParams:   config.DefaultPageParams,
		PrimaryBtn:   contactURL(),
		SecondaryBtn: "/resume",
		Profile:      *user,
		Description:  processBio(user.Bio),
	}

	return c.Render("index", props)
}

func resumeHandler(c *fiber.Ctx) error {
	logger.Info("building resume page")

	user, _ := github.User("marco-souza", "")
	pageParams := config.DefaultPageParams

	// set custom title
	pageParams.Title = fmt.Sprintf("Resume - %s", pageParams.Title)

	props := rootProps{
		PageParams:  pageParams,
		PrimaryBtn:  contactURL(),
		Profile:     *user,
		Description: processBio(user.Bio),
	}

	return c.Render("resume", props)
}

func blogHandler(c *fiber.Ctx) error {
	logger.Info("building blog page")

	user, _ := github.User("marco-souza", "")
	pageParams := config.DefaultPageParams

	// set custom title
	pageParams.Title = fmt.Sprintf("Blog - %s", pageParams.Title)

	props := rootProps{
		PageParams:  pageParams,
		PrimaryBtn:  contactURL(),
		Profile:     *user,
		Description: processBio(user.Bio),
	}

	return c.Render("blog", props)
}

func contactURL() string {
	q := url.Values{}
	q.Add("subject", "Hi Marco, Let's have a coffee")

	// mailto does not work with spaces as '+'
	contact := "mailto:marco@tremtec.com?" + strings.ReplaceAll(
		q.Encode(), "+", "%20",
	)

	logger.Info("contact Link generated", "contact", contact)

	return contact
}

var linkMap = map[string]string{
	"tremtec":  "https://tremtec.com",
	"podcodar": "https://podcodar.com",
	"mongodb":  "https://mongodb.com",
}

func processBio(text string) template.HTML {
	if text == "" {
		return ""
	}

	tagRegex := regexp.MustCompile(`@\w*`)

	result := tagRegex.ReplaceAllStringFunc(text, func(tag string) string {
		name := strings.TrimPrefix(tag, "@")
		link, ok := linkMap[name]
		if ok {
			return fmt.Sprintf(
				`<a class="text-pink-400" target="_blank" href="%s">@%s</a>`,
				link,
				name,
			)
		}
		return tag
	})

	return template.HTML(result)
}
