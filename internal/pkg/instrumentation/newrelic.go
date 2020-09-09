package instrumentation

import (
	"log"
	"net/http"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
)

type Config struct {
	Enabled bool
	AppName string
	License string
}

func MakeNewRelicApplication(conf Config) *newrelic.Application {
	app, err := newrelic.NewApplication(
		newrelic.ConfigEnabled(conf.Enabled),
		newrelic.ConfigAppName(conf.AppName),
		newrelic.ConfigLicense(conf.License),
		newrelic.ConfigDistributedTracerEnabled(true),
		func(c *newrelic.Config) {
			c.ErrorCollector.IgnoreStatusCodes = []int{http.StatusBadRequest, http.StatusNotFound, http.StatusUnprocessableEntity, http.StatusTooManyRequests, http.StatusConflict}
		},
	)
	if nil != err {
		log.Fatalf("error in init newrelic: %s", err)
	}

	// Wait for the application to connect.
	if err = app.WaitForConnection(5 * time.Second); nil != err {
		log.Fatalf("error in init newrelic: %s", err)
	}

	return app
}
