package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/tekumara/ghapp-dagger/pkg/app"

	"github.com/gregjones/httpcache"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/rcrowley/go-metrics"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	zerolog.DefaultContextLogger = &logger
	metricsRegistry := metrics.DefaultRegistry

	config, err := app.GhaConfig()
	if err != nil {
		log.Fatal(err)
	}

	cc, err := githubapp.NewDefaultCachingClientCreator(
		*config,
		githubapp.WithClientUserAgent("ghapp-dagger"),
		githubapp.WithClientTimeout(3*time.Second),
		githubapp.WithClientCaching(false, func() httpcache.Cache { return httpcache.NewMemoryCache() }),
		githubapp.WithClientMiddleware(
			githubapp.ClientMetrics(metricsRegistry),
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	checkSuiteHandler := &app.CheckSuiteEventHandler{
		ClientCreator: cc,
	}

	checkRunHandler := &app.CheckRunEventHandler{
		ClientCreator: cc,
		AppID: config.App.IntegrationID,
	}

	handler := githubapp.NewDefaultEventDispatcher(*config, checkSuiteHandler, checkRunHandler)

	http.Handle("/", handler)

	addr := fmt.Sprintf("%s:%d", "127.0.0.1", 8000)
	logger.Info().Msgf("Starting server on %s...", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
