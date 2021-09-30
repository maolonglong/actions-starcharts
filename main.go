package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/maolonglong/actions-starcharts/internal/action"
	"github.com/maolonglong/actions-starcharts/internal/chart"
	"github.com/maolonglong/actions-starcharts/internal/client"
	"github.com/spf13/cast"
)

func main() {
	owner, name := action.GetRepo()
	sha := action.GetSHA()
	token := action.GetInput("github_token")
	svgPath := action.GetInput("svg_path")
	commitMessage := action.GetInput("commit_message")

	starsChange := cast.ToInt(action.GetInput("stars_change"))
	if starsChange < 1 {
		starsChange = 1
	}

	client := client.New(context.Background(), token)

	cur, err := client.GetStarTotal(owner, name)
	if err != nil {
		log.Fatal("get stars total count failed: ", err.Error())
	}
	if cur == 0 {
		log.Println("not enough stars")
		os.Exit(0)
	}

	b, err := client.GetBlob(owner, name, sha, svgPath)
	if err != nil {
		log.Fatal("get blob failed: ", err.Error())
	}

	var old int
	fmt.Sscanf(b.Text, "<!-- stars: %d -->", &old)
	log.Printf("old stars: %v, cur stars: %v", old, cur)
	if abs(cur-old) < starsChange {
		os.Exit(0)
	}

	stars, err := client.GetStargazers(owner, name)
	if err != nil {
		log.Fatal("get stargazers failed: ", err.Error())
	}

	buf := new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("<!-- stars: %d -->\n", len(stars)))
	err = chart.WriteStarsChart(stars, buf)
	if err != nil {
		log.Fatal("write stargazers chart failed: ", err.Error())
	}

	err = client.CreateOrUpdate(owner, name, sha, svgPath, commitMessage, buf.Bytes())
	if err != nil {
		log.Fatal("update content failed: ", err.Error())
	}
}
