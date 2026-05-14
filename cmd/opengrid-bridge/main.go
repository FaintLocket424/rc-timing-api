package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/FaintLocket424/opengrid-bridge/internal/api"
	"github.com/FaintLocket424/opengrid-bridge/internal/manager"
	"github.com/FaintLocket424/opengrid-bridge/internal/storage/cache"
)

var (
	Version  = "dev"
	LogLevel = new(slog.LevelVar)
)

func InitLogger(useJSON bool) {
	LogLevel.Set(slog.LevelInfo)

	opts := &slog.HandlerOptions{
		Level: LogLevel,
	}

	var handler slog.Handler
	if useJSON {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	slog.SetDefault(slog.New(handler))
}

func main() {
	defaultPort := 4998

	if envPort := os.Getenv("PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			defaultPort = p
		} else {
			slog.Warn("Invalid PORT environment variable, falling back to default", "envPort", envPort, "default", defaultPort)
		}
	}

	port := flag.Int("port", defaultPort, "Port to listen on (overrides PORT env var)")
	host := flag.String("host", "0.0.0.0", "Host interface to bind to (e.g., 127.0.0.1 for localhost only)")
	jsonMode := flag.Bool("json", false, "Output logs in json format")
	flag.Parse()

	InitLogger(*jsonMode)

	if debugMode := Version == "dev" || strings.HasSuffix(Version, "debug"); debugMode {
		LogLevel.Set(slog.LevelDebug)
	}

	listenAddr := fmt.Sprintf("%s:%d", *host, *port)

	slog.Info("Starting OpenGrid Timing Bridge", "version", Version, "address", listenAddr)

	cache := cache.NewCache()
	manager := manager.NewManager(cache)
	router := api.SetupRouter(cache, manager)

	if err := router.Run(listenAddr); err != nil {
		slog.Error("An error occurred", "err", err)
		os.Exit(1)
	}
}
