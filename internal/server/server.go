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
	"github.com/marco-souza/marco.fly.dev/internal/lua"
	"github.com/marco-souza/marco.fly.dev/internal/server/routes"
	"github.com/marco-souza/marco.fly.dev/internal/telegram"
)

var logger = slog.With("service", "server")

type Server struct {
	IsProduction bool
	Done         *chan bool
	addr         string
	hostname     string
	port         string
	app          *fiber.App
}

func New(done *chan bool) *Server {
	conf := config.Load()
	hostname := conf.Hostname
	port := conf.Port
	addr := hostname + ":" + port

	engine := html.New(conf.Views, ".html")
	if conf.Env == "development" {
		engine.Debug(true)
		engine.Reload(true)
	}

	return &Server{
		Done:         done,
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

func (s *Server) Start() error {
	logger.Info("setting up routes")

	duration := 10 * time.Second
	if s.IsProduction {
		duration = 15 * time.Minute
	}

	s.app.Static("/static", "./static", fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		CacheDuration: duration,
		MaxAge:        3600,
	})

	routes.SetupRoutes(s.app)

	logger.Info("starting server ependencies")
	// order matters, as for now each service ask for its dependencies
	di.Injectables(
		config.Load,
		db.New,
		cache.New,
		discord.New,
		binance.New,
		telegram.New,
		lua.NewLuaService,
		cron.New,
	)

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

		if s.Done != nil {
			*s.Done <- true
		}

		return nil
	})

	return nil
}

func (s *Server) Run() error {
	// graceful shutdown on interrupt signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt) // register channel to interrupt signals
	go func() {
		<-shutdown // wait for shutdown signal
		logger.Info("shutting down server", "reason", "interrupt")
		di.Clean()
	}()

	// await for server to shutdown
	return s.app.Listen(s.addr)
}

func (s *Server) setupRoutes() {
}

func (s *Server) Stop() error {
	if err := s.app.Shutdown(); err != nil {
		logger.Warn("failed to shutdown server", "err", err)
	}

	return nil
}

func (s *Server) Test(req *http.Request, timeout ...int) (*http.Response, error) {
	return s.app.Test(req, timeout...)
}
