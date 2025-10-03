package main

import (
	"dashboard/internal/config"
	"dashboard/ui"
	"net/http"
	"os"

	"github.com/a-h/templ"
)

const (
	DashboardConfigPathEnvVar = "DASHBOARD_CONFIG_PATH"
)

func main() {
	cfg, err := loadConfiguration()
	if err != nil {
		panic(err)
	}

	services := ui.ServiceView(cfg.Services, cfg.Tags)
	layout := ui.Layout(cfg.Title, cfg.BackgroundColor, services)

	http.Handle("/", templ.Handler(layout))

	http.ListenAndServe(":8080", nil)
}

// loadConfiguration loads the configuration from the path specified in the DASHBOARD_CONFIG_PATH environment variable
func loadConfiguration() (*config.Config, error) {
	configPath := os.Getenv(DashboardConfigPathEnvVar)
	return config.LoadConfig(configPath)
}
