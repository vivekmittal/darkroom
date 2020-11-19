package config

import (
	newrelic "github.com/newrelic/go-agent"
)

func readNewRelicConfig() newrelic.Config {
	v := Viper()
	return newrelic.Config{
		Enabled: v.GetBool("newrelic.enabled"),
		AppName: v.GetString("newrelic.appName"),
		License: v.GetString("newrelic.license"),
	}
}
