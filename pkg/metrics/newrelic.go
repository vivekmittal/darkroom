package metrics

import (
	"github.com/gojek/darkroom/pkg/config"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func NewrelicApp() *newrelic.Application {
	newRelicApp, err := newrelic.NewApplication(
		newrelic.ConfigAppName(config.NewrelicConfig().AppName),
		newrelic.ConfigEnabled(config.NewrelicConfig().Enabled),
		newrelic.ConfigLicense(config.NewrelicConfig().License))
	if err != nil {
		panic(err)
	}

	return newRelicApp
}
