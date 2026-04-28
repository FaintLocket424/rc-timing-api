package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// example usage: `go run ./cmd/mirror -url "https://forcc.co.uk/live" -path "./internal/scraper/bbk/testdata/forcc_2026-04-10"

	targetURL := flag.String("url", "", "The base URL to mirror (e.g. https://forcc.co.uk/live)")
	outputPath := flag.String("path", ".", "The directory to save the mirrored files (defaults to current directory)")
	software := flag.String("software", "bbk", "The race timing software being scraped (defaults to bbk)")
	removeNonHTML := flag.Bool("remove-non-htm", false, "Remove files not ending in `.htm` after download")
	flag.Parse()

	if *targetURL == "" {
		fmt.Println("Error: the -url flag is required.")
		flag.Usage()
		os.Exit(1)
	}

	base := strings.TrimSuffix(*targetURL, "/")

	dir := time.Now().UTC().Format("2006-01-02_15-04-05")
	dateTimePath := filepath.Join(dir, *outputPath)

	var cmd *exec.Cmd

	switch *software {
	case "bbk":
		cmd = CreateMirrorCommandBBK(base, dateTimePath)
	default:
		log.Fatalf("Unable to parse software type: %s", *software)
	}

	log.Printf("Mirroring %s to %s...\n", base, *outputPath)

	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to run wget: %v", err)
	}

	if *removeNonHTML {
		if err := cleanupDirectory(*outputPath); err != nil {
			log.Fatalf("Cleanup failed: %v", err)
		}
	}

	log.Println("Mirroring complete.")
}

func cleanupDirectory(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		keep := strings.HasSuffix(info.Name(), ".htm") || strings.HasSuffix(info.Name(), ".html")

		if !keep {
			return os.Remove(path)
		}

		return nil
	})
}
