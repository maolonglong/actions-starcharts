package main

import (
	"os"
	"strings"

	"github.com/sethvargo/go-githubactions"
)

// TODO: Use https://github.com/sethvargo/go-githubactions/pull/46 instead
func getRepo() (string, string) {
	a := strings.SplitN(os.Getenv("GITHUB_REPOSITORY"), "/", 2)
	if len(a) != 2 || a[0] == "" || a[1] == "" {
		githubactions.Fatalf("failed to get repo")
	}
	return a[0], a[1]
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
