package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/sethvargo/go-githubactions"
)

func main() {
	ghCtx, err := githubactions.Context()
	if err != nil {
		githubactions.Fatalf("get github context failed")
	}
	sha := ghCtx.SHA
	owner, name := getRepo()

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

	starsChange, err := strconv.Atoi(githubactions.GetInput("stars_change"))
	if err != nil || starsChange < 1 {
		starsChange = 1
	}

	targetOwner, targetName := owner, name
	repo := githubactions.GetInput("repo")
	if repo != "" {
		a := strings.SplitN(repo, "/", 2)
		if len(a) != 2 {
			githubactions.Fatalf("invalid repo: %v", repo)
		}
		targetOwner, targetName = a[0], a[1]
	}

	ctx := context.TODO()
	cli := newClient(token)

	cur, err := cli.getStarsCount(ctx, targetOwner, targetName)
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
		preContent, err := base64.StdEncoding.DecodeString(*b.Content)
		if err != nil {
			githubactions.Fatalf("failed to decode base64 string: %v", err)
		}

		var old int
		fmt.Sscanf(string(preContent), "<!-- stars: %d -->", &old)
		githubactions.Infof("stars_old=%d stars_cur=%d", old, cur)
		if abs(cur-old) < starsChange {
			os.Exit(0)
		}
	}

	stars, err := cli.getStargazers(ctx, targetOwner, targetName)
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
