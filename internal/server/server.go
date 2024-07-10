package server

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"

	"github.com/marco-souza/marco.fly.dev/internal/cache"
	"github.com/marco-souza/marco.fly.dev/internal/config"
	"github.com/marco-souza/marco.fly.dev/internal/cron"
	"github.com/marco-souza/marco.fly.dev/internal/db"
	"github.com/marco-souza/marco.fly.dev/internal/discord"
	"github.com/marco-souza/marco.fly.dev/internal/server/routes"
)

type server struct {
	addr     string
	hostname string
	port     string
	app      *fiber.App
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
		addr:     addr,
		port:     port,
		hostname: hostname,
		app: fiber.New(fiber.Config{
			Views:       engine,
			ViewsLayout: "layouts/main",
		}),
	}
}

func (s *server) Start(done *chan bool) {
	log.Println("setting up routes")
	s.setupRoutes()

	// TODO: seed sqlc db

	startup := func() error {
		log.Println("starting server dependencies")

		if err := db.Init(config.Load().SqliteUrl); err != nil {
			return err
		}

		if err := cron.Start(); err != nil {
			log.Println("[warn] failed to start cron:", err)
		}

		if err := discord.DiscordService.Open(); err != nil {
			log.Println("[warn] failed to start discord:", err)
		}

		if err := cache.SetStorage(cache.NewMemCache()); err != nil {
			log.Println("[warn] failed to start cache:", err)
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
			log.Println("listening on: ", scheme+"://"+listenData.Host+":"+listenData.Port)

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
		log.Fatal(err)
	}
}

func (s *server) setupRoutes() {
	log.Println("setup static resources")
	s.app.Static("/static", "./static", fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		CacheDuration: 10 * time.Second,
		MaxAge:        3600,
	})

	routes.SetupRoutes(s.app)
}

func (s *server) Shutdown() {
	log.Println("shutting down server dependencies")

	cron.Stop()
	discord.DiscordService.Close()

	if err := s.app.Shutdown(); err != nil {
		log.Println("failed to shutdown server", err)
	}

	if err := db.Close(); err != nil {
		log.Println("failed to shutdown db", err)
	}

	log.Println("bye bye!")
}

func (s *server) Test(req *http.Request, timeout ...int) (*http.Response, error) {
	return s.app.Test(req, timeout...)
}
