package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

func CreateMirrorCommandBBK(base, outputPath string) *exec.Cmd {
	urls := []string{
		base + "/",
		base + "/liveraceres.htm",
		base + "/liveresults.htm",
		base + "/liveschedule.htm",
		base + "/livecompets.htm",
	}

	args := []string{
		"--mirror",
		"--convert-links",
		"--page-requisites",
		"--no-parent",
		"--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
		"--directory-prefix=" + outputPath,
		"--input-file=-",
	}

	cmd := exec.Command("wget", args...)

	cmd.Stdin = bytes.NewBufferString(strings.Join(urls, "\n"))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}
