package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"sort"
	"time"

	"github.com/google/go-github/v39/github"
	"github.com/sethvargo/go-githubactions"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type client struct {
	v3 *github.Client
	v4 *githubv4.Client
}

func newClient(token string) *client {
	ctx := context.TODO()

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(ctx, src)

	return &client{
		v3: github.NewClient(httpClient),
		v4: githubv4.NewClient(httpClient),
	}
}

func (c *client) getStarsCount(ctx context.Context, owner, repo string) (int, error) {
	r, _, err := c.v3.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return 0, err
	}
	if r.StargazersCount == nil {
		return 0, nil
	}
	return *r.StargazersCount, nil
}

func (c *client) getBlob(ctx context.Context, owner, repo, sha, path string) (*github.Blob, error) {
	tree, _, err := c.v3.Git.GetTree(ctx, owner, repo, sha, true)
	if err != nil {
		return nil, err
	}
	for _, ent := range tree.Entries {
		if ent.Path != nil && *ent.Path == path {
			blob, _, err := c.v3.Git.GetBlob(ctx, owner, repo, ent.GetSHA())
			if err != nil {
				return nil, err
			}
			return blob, nil
		}
	}
	return nil, nil
}

type stargazer struct {
	StarredAt time.Time
}

type getStarsQuery struct {
	Repository struct {
		Stargazers struct {
			Edges    []stargazer
			PageInfo struct {
				EndCursor   githubv4.String
				HasNextPage bool
			}
		} `graphql:"stargazers(first: 100, after: $after)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
	RateLimit struct {
		Remaining int
	}
}

func (c *client) getStargazers(ctx context.Context, owner, repo string) ([]stargazer, error) {
	var (
		q      getStarsQuery
		params = map[string]any{
			"owner": githubv4.String(owner),
			"name":  githubv4.String(repo),
			"after": (*githubv4.String)(nil),
		}
		stars []stargazer
		last  time.Time
	)

	githubactions.Group(fmt.Sprintf("getStargazers(%q, %q)", owner, repo))

	for {
		err := c.v4.Query(ctx, &q, params)
		if err != nil {
			return nil, err
		}

		stars = append(stars, q.Repository.Stargazers.Edges...)
		if last.IsZero() || time.Since(last) > 2*time.Second {
			githubactions.Infof("fetching: count=%d ratelimit_remaining=%d",
				len(stars), q.RateLimit.Remaining)
			last = time.Now()
		}

		if !q.Repository.Stargazers.PageInfo.HasNextPage {
			break
		}
		params["after"] = githubv4.NewString(q.Repository.Stargazers.PageInfo.EndCursor)
	}

	githubactions.Infof("completed: count=%d ratelimit_remaining=%d",
		len(stars), q.RateLimit.Remaining)

	githubactions.EndGroup()

	sort.Slice(stars, func(i, j int) bool {
		return stars[i].StarredAt.Before(stars[j].StarredAt)
	})
	return stars, nil
}

func (c *client) createOrUpdate(ctx context.Context, owner, repo, sha, path, message string, blob *github.Blob, content []byte) error {
	if blob != nil {
		preContent, err := base64.StdEncoding.DecodeString(*blob.Content)
		if err != nil {
			return err
		}

		if bytes.Equal(preContent, content) {
			return nil
		}
	}

	opt := &github.RepositoryContentFileOptions{
		Message: &message,
		Content: content,
	}
	if blob != nil {
		opt.SHA = blob.SHA
	}
	_, _, err := c.v3.Repositories.CreateFile(ctx, owner, repo, path, opt)
	return err
}
