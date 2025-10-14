package main

import (
	"dashboard/internal/api"
	"dashboard/internal/config"
	"net/http"
	"os"
)

const (
	DashboardConfigPathEnvVar = "DASHBOARD_CONFIG_PATH"
)

func main() {
	cfg, err := loadConfiguration()
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/config", api.ConfigHandler(cfg))

	// Serve Static Frontend
	mux.Handle("/", http.FileServer(http.Dir("./frontend/dist")))

	http.ListenAndServe("[::]:8080", api.CORSMiddleware(mux))
}

// loadConfiguration loads the configuration from the path specified in the DASHBOARD_CONFIG_PATH environment variable
func loadConfiguration() (*config.Config, error) {
	configPath := os.Getenv(DashboardConfigPathEnvVar)
	return config.LoadConfig(configPath)
}
