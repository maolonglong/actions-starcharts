package action

import (
	"os"
	"strings"
)

func GetSHA() string {
	return os.Getenv("GITHUB_SHA")
}

func GetRepo() (owner string, repo string) {
	a := strings.SplitN(os.Getenv("GITHUB_REPOSITORY"), "/", 2)
	return a[0], a[1]
}

func GetInput(name string) string {
	return strings.TrimSpace(os.Getenv("INPUT_" + strings.ToUpper(name)))
}
