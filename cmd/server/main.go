package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/FaintLocket424/rc-timing-api/internal/api"
	"github.com/FaintLocket424/rc-timing-api/internal/manager"
	"github.com/FaintLocket424/rc-timing-api/internal/storage/cache"
	"github.com/gin-gonic/gin"
)

var (
	version         = "dev"
	GinMode  string = "debug"
	LogLevel        = new(slog.LevelVar)
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
	if GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	debugMode := version == "dev"
	jsonMode := flag.Bool("json", false, "output logs in json format")
	flag.Parse()

	InitLogger(*jsonMode)

	if debugMode {
		LogLevel.Set(slog.LevelDebug)
	}

	slog.Info("Starting RC Timing API", "version", version)

	cache := cache.NewCache()
	manager := manager.NewManager(cache)
	router := api.SetupRouter(cache, manager)

	if err := router.Run("0.0.0.0:8080"); err != nil {
		slog.Error("An error occurred", "err", err)
		return
	}
}
