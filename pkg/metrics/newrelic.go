package metrics

import (
	"github.com/gojek/darkroom/pkg/config"
	newrelic "github.com/newrelic/go-agent"
)

func NewrelicApp() newrelic.Application {
	newRelicApp, err := newrelic.NewApplication(config.NewrelicConfig())
	if err != nil {
		panic(err)
	}

	return newRelicApp
}
