package client

import (
	"time"

	"github.com/shurcooL/githubv4"
)

type Stargazer struct {
	StarredAt time.Time
}

/*
repository(name: $name, owner: $owner) {
	stargazers(first: 100, after: $after) {
		edges {
			starredAt
		}
		pageInfo {
			endCursor
			hasNextPage
		}
	}
}
*/
type getStarsQuery struct {
	Repository struct {
		Stargazers struct {
			Edges    []Stargazer
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

type Blob struct {
	Oid  string
	Text string
}

/*
{
  repository(owner: $owner, name: $name) {
    object(expression: $expression) {
      ... on Blob {
        oid
		text
      }
    }
  }
}
*/
type getFileSHAQuery struct {
	Repository struct {
		Object struct {
			Blob `graphql:"... on Blob"`
		} `graphql:"object(expression: $expression)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}
