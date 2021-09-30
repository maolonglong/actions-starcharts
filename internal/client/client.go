package client

import (
	"context"
	"net/http"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type Client struct {
	ctx context.Context
	g   *githubv4.Client

	// TODO: remove httpClient and github api v3
	httpClient *http.Client
}

func New(ctx context.Context, token string) *Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(ctx, src)
	return &Client{
		ctx:        ctx,
		httpClient: httpClient,
		g:          githubv4.NewClient(httpClient),
	}
}
