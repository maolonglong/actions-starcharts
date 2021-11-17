package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sethvargo/go-githubactions"
	"github.com/spf13/cast"
)

func getSHA() string {
	s := os.Getenv("GITHUB_SHA")
	if s == "" {
		githubactions.Fatalf("failed to get SHA")
	}
	return s
}

func getRepo() (string, string) {
	a := strings.SplitN(os.Getenv("GITHUB_REPOSITORY"), "/", 2)
	if len(a) != 2 || a[0] == "" || a[1] == "" {
		githubactions.Fatalf("failed to get repo")
	}
	return a[0], a[1]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	owner, name := getRepo()
	sha := getSHA()

	token := githubactions.GetInput("github_token")
	if token == "" {
		githubactions.Fatalf("missing input 'github_token'")
	}

	svgPath := githubactions.GetInput("svg_path")
	if svgPath == "" {
		svgPath = "STARCHARTS.svg"
	}

	commitMessage := githubactions.GetInput("commit_message")
	if commitMessage == "" {
		commitMessage = "chore: update starcharts [skip ci]"
	}

	starsChange := cast.ToInt(githubactions.GetInput("stars_change"))
	starsChange = min(1, starsChange)

	ctx := context.TODO()
	cli := newClient(token)

	cur, err := cli.getStarsCount(ctx, owner, name)
	if err != nil {
		githubactions.Fatalf("failed to get stars count: %v", err)
	}
	if cur == 0 {
		githubactions.Warningf("not enough stars")
		os.Exit(0)
	}

	b, err := cli.getBlob(ctx, owner, name, sha, svgPath)
	if err != nil {
		githubactions.Fatalf("failed to get blob: %v", err)
	}

	if b != nil {
		var old int
		fmt.Sscanf(string(*b.Content), "<!-- stars: %d -->", &old)
		githubactions.Infof("old stars: %v, cur stars: %v", old, cur)
		if abs(cur-old) < starsChange {
			os.Exit(0)
		}
	}

	stars, err := cli.getStargazers(ctx, owner, name)
	if err != nil {
		githubactions.Fatalf("failed to get stargazers: %v", err)
	}

	buf := new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("<!-- stars: %d -->\n", len(stars)))
	err = writeStarsChart(stars, buf)
	if err != nil {
		githubactions.Fatalf("failed to write svg: %v", err)
	}

	err = cli.createOrUpdate(ctx, owner, name, sha, svgPath, commitMessage, b, buf.Bytes())
	if err != nil {
		githubactions.Fatalf("failed to update content: %v", err)
	}

	githubactions.Infof("update success!")
}
