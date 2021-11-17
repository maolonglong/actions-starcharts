package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/sethvargo/go-githubactions"
	"github.com/spf13/cast"
)

func getSHA() string {
	s := os.Getenv("GITHUB_SHA")
	if s == "" {
		githubactions.Fatalf("failed to get SHA\n")
	}
	return s
}

func getRepo() (string, string) {
	a := strings.SplitN(os.Getenv("GITHUB_REPOSITORY"), "/", 2)
	if len(a) != 2 || a[0] == "" || a[1] == "" {
		githubactions.Fatalf("failed to get repo\n")
	}
	return a[0], a[1]
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
		githubactions.Fatalf("missing input 'github_token'\n")
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
	if starsChange < 1 {
		starsChange = 1
	}

	targetOwner, targetName := owner, name
	repo := githubactions.GetInput("repo")
	if repo != "" {
		a := strings.SplitN(repo, "/", 2)
		if len(a) != 2 {
			githubactions.Fatalf("invalid repo: %v\n", repo)
		}
		targetOwner, targetName = a[0], a[1]
	}

	ctx := context.TODO()
	cli := newClient(token)

	cur, err := cli.getStarsCount(ctx, targetOwner, targetName)
	if err != nil {
		githubactions.Fatalf("failed to get stars count: %v\n", err)
	}
	if cur == 0 {
		githubactions.Warningf("not enough stars\n")
		os.Exit(0)
	}

	b, err := cli.getBlob(ctx, owner, name, sha, svgPath)
	if err != nil {
		githubactions.Fatalf("failed to get blob: %v\n", err)
	}

	if b != nil {
		preContent, err := base64.StdEncoding.DecodeString(*b.Content)
		if err != nil {
			githubactions.Fatalf("failed to decode base64 string: %v\n", err)
		}

		var old int
		fmt.Sscanf(string(preContent), "<!-- stars: %d -->", &old)
		githubactions.Infof("old_stars=%d cur_stars=%d\n", old, cur)
		if abs(cur-old) < starsChange {
			os.Exit(0)
		}
	}

	stars, err := cli.getStargazers(ctx, targetOwner, targetName)
	if err != nil {
		githubactions.Fatalf("failed to get stargazers: %v\n", err)
	}

	buf := new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("<!-- stars: %d -->\n", len(stars)))
	err = writeStarsChart(stars, buf)
	if err != nil {
		githubactions.Fatalf("failed to write svg: %v\n", err)
	}

	err = cli.createOrUpdate(ctx, owner, name, sha, svgPath, commitMessage, b, buf.Bytes())
	if err != nil {
		githubactions.Fatalf("failed to update content: %v\n", err)
	}

	githubactions.Infof("update success!\n")
}
