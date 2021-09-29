package client

import (
	"context"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

const defaultPerPage = 1000

type Client struct {
	g   *github.Client
	ctx context.Context
}

func New(ctx context.Context, token string) *Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return &Client{g: github.NewClient(tc), ctx: ctx}
}
