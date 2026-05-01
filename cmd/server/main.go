package main

import (
	"flag"
	"log/slog"
	"os"
	"strings"

	"github.com/FaintLocket424/rc-timing-api/internal/api"
	"github.com/FaintLocket424/rc-timing-api/internal/manager"
	"github.com/FaintLocket424/rc-timing-api/internal/storage/cache"
)

var (
	Version     = "dev"
	DefaultPort = "8080"
	LogLevel    = new(slog.LevelVar)
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
	jsonMode := flag.Bool("json", false, "output logs in json format")
	flag.Parse()
	InitLogger(*jsonMode)

	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}

	debugMode := Version == "dev" || strings.HasSuffix(Version, "debug")

	if debugMode {
		LogLevel.Set(slog.LevelDebug)
	}

	slog.Info("Starting RC Timing API", "version", Version, "port", port)

	cache := cache.NewCache()
	manager := manager.NewManager(cache)
	router := api.SetupRouter(cache, manager)

	if err := router.Run("0.0.0.0:" + port); err != nil {
		slog.Error("An error occurred", "err", err)
		return
	}
}
