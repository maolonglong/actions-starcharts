package main

import (
	"bytes"
	"context"
	"log"

	"github.com/maolonglong/actions-starcharts/internal/action"
	"github.com/maolonglong/actions-starcharts/internal/chart"
	"github.com/maolonglong/actions-starcharts/internal/client"
)

func main() {
	owner, repo := action.GetRepo()
	sha := action.GetSHA()
	token := action.GetInput("github_token")
	svgPath := action.GetInput("svg_path")
	commitMessage := action.GetInput("commit_message")

	client := client.New(context.Background(), token)
	stars, err := client.GetStargazers(owner, repo)
	if err != nil {
		log.Fatal("get stargazers failed: ", err.Error())
	}

	buf := new(bytes.Buffer)
	err = chart.WriteStarsChart(stars, buf)
	if err != nil {
		log.Fatal("write stargazers chart failed: ", err.Error())
	}

	err = client.CreateOrUpdate(owner, repo, sha, svgPath, commitMessage, buf.Bytes())
	if err != nil {
		log.Fatal("update content failed: ", err.Error())
	}
}
