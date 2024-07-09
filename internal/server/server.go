package server

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"

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

var conf = config.Load()

func New() *server {
	hostname := conf.Hostname
	port := conf.Port
	addr := fmt.Sprintf("%s:%s", hostname, port)

	engine := html.New("./views", ".html")
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

func (s *server) Start() {
	fmt.Println("setting up routes...")
	s.setupRoutes()

	// TODO: seed sqlc db

	startup := func() error {
		fmt.Println("starting services...")

		if err := db.Init(conf.SqliteUrl); err != nil {
			return err
		}

		if err := cron.Start(); err != nil {
			return err
		}

		if err := discord.DiscordService.Open(); err != nil {
			return err
		}

		return s.app.Listen(s.addr)
	}

	teardown := func() {
		fmt.Println("shutting down services...")
		db.Close() // TODO: deprecate

		cron.Stop()
		discord.DiscordService.Close()

		if err := s.app.Shutdown(); err != nil {
			log.Fatal(err)
		}

		if err := db.Close(); err != nil {
			log.Fatalf("error closing db: %e", err)
		}
	}

	// graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt) // register channel to interrupt signals
	go func() {
		<-shutdown // wait for shutdown signal
		teardown()
	}()

	// await for server to shutdown
	if err := startup(); err != nil {
		teardown()
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
