package client

import (
	"log"
	"sort"

	"github.com/shurcooL/githubv4"
)

func (c *Client) GetStargazers(owner, name string) ([]Stargazer, error) {
	var q getStarsQuery

	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(name),
		"after": (*githubv4.String)(nil),
	}

	var stars []Stargazer
	for {
		err := c.g.Query(c.ctx, &q, variables)
		if err != nil {
			return nil, err
		}
		stars = append(stars, q.Repository.Stargazers.Edges...)
		log.Printf("get stargazers: completed=%v, ratelimit_remaining=%v", len(stars), q.RateLimit.Remaining)
		if !q.Repository.Stargazers.PageInfo.HasNextPage {
			break
		}
		variables["after"] = githubv4.NewString(q.Repository.Stargazers.PageInfo.EndCursor)
	}

	sort.Slice(stars, func(i, j int) bool {
		return stars[i].StarredAt.Before(stars[j].StarredAt)
	})
	return stars, nil
}
