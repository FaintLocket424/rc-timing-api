package mirror

import (
	"bytes"
	"flag"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	targetURL := flag.String("url", "", "The base URL to mirror (e.g. https://forcc.co.uk/live)")
	outputPath := flag.String("path", ".", "The directory to save the mirrored files (defaults to current directory)")
	flag.Parse()

	if *targetURL == "" && flag.NArg() > 0 {
		*targetURL = flag.Arg(0)
	}

	if *outputPath == "." && flag.NArg() > 1 {
		*outputPath = flag.Arg(1)
	}

	if *targetURL == "" {
		log.Fatalln("Error: Please provide a target URL.\nUsage: go run ./cmd/mirror -url \"https://forcc.co.uk/live\"")
	}

	base := strings.TrimSuffix(*targetURL, "/")

	urls := []string{
		base + "/",
		base + "/liveraceres.htm",
		base + "liveresults.htm",
		base + "/liveschedule.htm",
		base + "/livecompets.htm",
	}

	args := []string{
		"--mirror",
		"--convert-links",
		"--page-requisites",
		"--no-parent",
		"--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
		"--directory-prefix=" + *outputPath,
		"--input-file=-",
	}

	cmd := exec.Command("wget", args...)

	cmd.Stdin = bytes.NewBufferString(strings.Join(urls, "\n"))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("Mirroring %s to %s...\n", base, *outputPath)

	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to run wget: %v", err)
	}
	log.Println("Mirroring complete.")
}
