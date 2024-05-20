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
	"github.com/marco-souza/marco.fly.dev/internal/models"
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

	if conf.Env == "development" {
		fmt.Println("seeding db...")
		models.Seed()
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt) // register channel to interrupt signals
	teardown := func() {
		<-shutdown
		fmt.Println("shutting down server...")
		s.app.Shutdown()
	}

	go teardown() // start listening for interrupt signals

	// await for server to shutdown
	if err := s.app.Listen(s.addr); err != nil {
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
