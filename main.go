package main

import (
	"flag"
	"fmt"
	"github.com/ben-burwood/dashboard/internal/config"
	"github.com/ben-burwood/dashboard/internal/render"
	"net/http"
	"os"
	"path/filepath"
)

const (
	DashboardConfigPathEnvVar = "DASHBOARD_CONFIG_PATH"

	// DefaultBuildOutputPath is where the static export is written when --build is used without --out
	DefaultBuildOutputPath = "build/index.html"
)

func main() {
	build := flag.Bool("build", false, "render the dashboard to a static HTML file and exit (instead of serving)")
	out := flag.String("out", DefaultBuildOutputPath, "output path for --build")
	configPath := flag.String("config", "", "config file or directory (overrides $"+DashboardConfigPathEnvVar+")")
	flag.Parse()

	cfg, err := loadConfiguration(*configPath)
	if err != nil {
		panic(err)
	}

	page, err := render.Render(cfg)
	if err != nil {
		panic(err)
	}

	if *build {
		if err := writeBuild(*out, page); err != nil {
			panic(err)
		}
		fmt.Printf("Wrote %s\n", *out)
		return
	}

	serve(page)
}

// serve renders the dashboard once at startup and serves it over HTTP on :8080
func serve(page []byte) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(page)
	})
	http.ListenAndServe("[::]:8080", mux)
}

// writeBuild writes the rendered page to path, creating parent directories as needed
func writeBuild(path string, page []byte) error {
	if dir := filepath.Dir(path); dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	return os.WriteFile(path, page, 0o644)
}

// loadConfiguration resolves the config path from the -config flag, falling back
// to the DASHBOARD_CONFIG_PATH environment variable. When neither is set, an empty
// path is passed and config.LoadConfig applies its own default (config/config.yml).
func loadConfiguration(configPath string) (*config.Config, error) {
	if configPath == "" {
		configPath = os.Getenv(DashboardConfigPathEnvVar)
	}
	return config.LoadConfig(configPath)
}
