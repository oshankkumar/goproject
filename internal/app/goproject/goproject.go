package goproject

import (
	"example.com/example/goproject/http/handler"
	"example.com/example/goproject/http/middleware"
	"example.com/example/goproject/http/router"
	"example.com/example/goproject/http/server"
	"example.com/example/goproject/internal/pkg/config"
	"example.com/example/goproject/internal/pkg/instrumentation"
	"example.com/example/goproject/pkg/i18n"
	"example.com/example/goproject/pkg/logger"
	"example.com/example/goproject/pkg/signal"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"gopkg.in/alexcesaro/statsd.v2"
)

type GoProjectAPP struct {
	apiServer    *server.Server
	Config       *config.Config
	newRelicApp  *newrelic.Application
	statsDClient *statsd.Client
	Tr           *i18n.Translator
	Deps         *GoProjectDeps
}

func NewAPP() *GoProjectAPP {
	return &GoProjectAPP{Config: config.New()}
}

// Initialise your application
func (g *GoProjectAPP) Init() error {
	logger.Init(logger.Config{
		Format: logger.Format(g.Config.LogFormat),
		Level:  g.Config.LogLevel,
	})

	g.newRelicApp = instrumentation.MakeNewRelicApplication(instrumentation.Config{
		Enabled: g.Config.NewRelicEnabled,
		AppName: g.Config.NewrelicAppName,
		License: g.Config.NewRelicKey,
	})

	g.Deps = NewGoProjectDeps(g.Config)

	g.Tr = i18n.MustTranslator(g.Config.TranslationPath)

	return nil
}

func (g *GoProjectAPP) Start() error {
	// Add Routers here
	routers := g.createRoutes()

	stop := signal.SetupHandler()
	errCh := make(chan error)

	g.apiServer = server.New(g.Config.Address, g.Tr)
	g.apiServer.InitRouter(routers...)
	// Add Global Middlewares here
	g.apiServer.Use(
		middleware.WithRecover,
		middleware.NewRelic(g.newRelicApp, "/metrics", "/docs", "/ping"),
		middleware.StatsD(g.statsDClient),
		middleware.WithLogger(),
	)

	go func() { errCh <- g.apiServer.Run() }()

	logrus.WithField("address", g.Config.Address).Info("started api server")

	select {
	case <-stop:
		return nil
	case err := <-errCh:
		return err
	}
}

func (g *GoProjectAPP) createRoutes() []router.Router {
	return []router.Router{
		handler.Ping(),
	}
}

// Close your application
func (g *GoProjectAPP) Close() error {
	return g.apiServer.Stop(g.Config.ShutDownTimeout)
}
