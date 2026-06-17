/*
 * OpenGrid Timing Bridge
 * Copyright (C) 2026 OpenGrid RC
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

// Package main is the entry point for the rc-timing-api web server.
//
// It creates a cache from the [internal/storage/cache] package.
// It creates a supervisor from the [internal/storage/tracking] package.
// It then uses the cache and supervisor to create a router in the [internal/api] package.
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/FaintLocket424/opengrid-bridge/internal/api"
	"github.com/FaintLocket424/opengrid-bridge/internal/scraper"
	"github.com/FaintLocket424/opengrid-bridge/internal/storage/cache"
	"github.com/FaintLocket424/opengrid-bridge/internal/tracking"
)

var (
	Version  = "dev"
	Commit   = "unknown"
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
	port := flag.Int("port", 4998, "Port to listen on (overrides PORT env var)")
	host := flag.String("host", "0.0.0.0", "Host interface to bind to (e.g., 127.0.0.1 for localhost only)")
	jsonMode := flag.Bool("json", false, "Output logs in json format")
	flag.Parse()

	InitLogger(*jsonMode)

	if debugMode := Version == "dev" || strings.HasSuffix(Version, "debug"); debugMode {
		LogLevel.Set(slog.LevelDebug)
	}

	listenAddr := fmt.Sprintf("%s:%d", *host, *port)

	slog.Info("Starting OpenGrid Timing Bridge", "version", Version, "commit", Commit, "address", listenAddr)

	cache := cache.NewCache()
	scraperFactory := scraper.NewFactory(Version, Commit)
	tracker := tracking.NewSupervisor(cache, scraperFactory)
	router := api.SetupRouter(cache, tracker, Version, Commit)

	if err := router.Run(listenAddr); err != nil {
		slog.Error("An error occurred", "err", err)
		os.Exit(1)
	}
}
