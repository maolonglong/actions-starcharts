package client

import (
	"sort"
	"sync"

	"github.com/google/go-github/v39/github"
	"golang.org/x/sync/errgroup"
)

func (c *Client) GetStargazers(owner, repo string) ([]*github.Stargazer, error) {
	r, _, err := c.g.Repositories.Get(c.ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	stars := make([]*github.Stargazer, 0, r.GetStargazersCount())
	var mu sync.Mutex
	var g errgroup.Group
	sem := make(chan struct{}, 4)

	lastPage := getLastPage(r.GetStargazersCount(), defaultPerPage)
	for page := 1; page <= lastPage; page++ {
		sem <- struct{}{}
		page := page
		g.Go(func() error {
			defer func() { <-sem }()
			opt := &github.ListOptions{
				Page:    page,
				PerPage: defaultPerPage,
			}
			result, _, err := c.g.Activity.ListStargazers(c.ctx, owner, repo, opt)
			if err != nil {
				return err
			}
			mu.Lock()
			stars = append(stars, result...)
			mu.Unlock()
			return nil
		})
	}

	err = g.Wait()
	if err != nil {
		return nil, err
	}

	sort.Slice(stars, func(i, j int) bool {
		return stars[i].StarredAt.Before(stars[j].StarredAt.Time)
	})
	return stars, nil
}

func getLastPage(total, perPage int) int {
	return (total + perPage - 1) / perPage
}
