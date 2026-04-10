package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	targetURL := flag.String("url", "", "The base URL to mirror (e.g. https://forcc.co.uk/live)")
	outputPath := flag.String("path", ".", "The directory to save the mirrored files (defaults to current directory)")
	software := flag.String("software", "bbk", "The race timing software being scraped (defaults to bbk)")
	flag.Parse()

	if *targetURL == "" {
		fmt.Println("Error: the -url flag is required.")
		flag.Usage()
		os.Exit(1)
	}

	base := strings.TrimSuffix(*targetURL, "/")

	var cmd *exec.Cmd

	switch *software {
	case "bbk":
		cmd = CreateMirrorCommandBBK(base, *outputPath)
	default:
		log.Fatalf("Unable to parse software type: %s", *software)
	}

	log.Printf("Mirroring %s to %s...\n", base, *outputPath)

	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to run wget: %v", err)
	}

	log.Println("Mirroring complete.")
}
