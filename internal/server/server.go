package server

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"

	"github.com/marco-souza/marco.fly.dev/internal/binance"
	"github.com/marco-souza/marco.fly.dev/internal/cache"
	"github.com/marco-souza/marco.fly.dev/internal/config"
	"github.com/marco-souza/marco.fly.dev/internal/cron"
	"github.com/marco-souza/marco.fly.dev/internal/db"
	"github.com/marco-souza/marco.fly.dev/internal/di"
	"github.com/marco-souza/marco.fly.dev/internal/discord"
	"github.com/marco-souza/marco.fly.dev/internal/server/routes"
	"github.com/marco-souza/marco.fly.dev/internal/telegram"
)

var logger = slog.With("service", "server")

type server struct {
	IsProduction bool
	addr         string
	hostname     string
	port         string
	app          *fiber.App
}

func New() *server {
	conf := config.Load()
	hostname := conf.Hostname
	port := conf.Port
	addr := hostname + ":" + port

	engine := html.New(conf.Views, ".html")
	if conf.Env == "development" {
		engine.Debug(true)
		engine.Reload(true)
	}

	return &server{
		IsProduction: conf.Env == "production",
		addr:         addr,
		port:         port,
		hostname:     hostname,
		app: fiber.New(fiber.Config{
			Views:       engine,
			ViewsLayout: "layouts/main",
		}),
	}
}

func (s *server) Start(done *chan bool) {
	logger.Info("setting up routes")
	s.setupRoutes()

	startup := func() error {
		if err := s.setupServices(); err != nil {
			return err
		}

		// listen for server events
		s.app.Hooks().OnListen(func(listenData fiber.ListenData) error {
			if fiber.IsChild() {
				return nil
			}
			scheme := "http"
			if listenData.TLS {
				scheme = "https"
			}
			url := scheme + "://" + listenData.Host + ":" + listenData.Port
			logger.Info("listening on " + url)

			if done != nil {
				*done <- true
			}

			return nil
		})

		return s.app.Listen(s.addr)
	}

	// graceful shutdown on interrupt signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt) // register channel to interrupt signals
	go func() {
		<-shutdown // wait for shutdown signal
		s.Shutdown()
	}()

	// await for server to shutdown
	if err := startup(); err != nil {
		s.Shutdown()
		logger.Error("server failed", "err", err)
	}
}

func (s *server) setupRoutes() {
	duration := 10 * time.Second
	if s.IsProduction {
		duration = 15 * time.Minute
	}

	logger.Info("setup static resources")
	s.app.Static("/static", "./static", fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		CacheDuration: duration,
		MaxAge:        3600,
	})

	routes.SetupRoutes(s.app)
}

func (s *server) setupServices() error {
	logger.Info("starting server dependencies")

	// order matters, as for now each service ask for its dependencies
	di.Injectable(config.Load())
	di.Injectable(db.New())
	di.Injectable(cache.New())
	di.Injectable(cron.New())
	di.Injectable(discord.New())
	di.Injectable(binance.New())
	di.Injectable(telegram.New())

	return nil
}

func (s *server) Shutdown() {
	di.Clean()

	if err := s.app.Shutdown(); err != nil {
		logger.Warn("failed to shutdown server", "err", err)
	}

	logger.Info("bye!")
}

func (s *server) Test(req *http.Request, timeout ...int) (*http.Response, error) {
	return s.app.Test(req, timeout...)
}
