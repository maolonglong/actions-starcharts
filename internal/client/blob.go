package client

import (
	"fmt"

	"github.com/google/go-github/v39/github"
	"github.com/shurcooL/githubv4"
)

func (c *Client) CreateOrUpdate(owner, repo, sha, path, message string, content []byte) error {
	b, err := c.GetBlob(owner, repo, sha, path)
	if err != nil {
		return err
	}

	if b.Text == string(content) {
		return nil
	}

	// TODO: replace with graphql API
	opt := &github.RepositoryContentFileOptions{
		Message: &message,
		Content: content,
	}
	if b.Oid != "" {
		opt.SHA = &b.Oid
	}
	v3 := github.NewClient(c.httpClient)
	_, _, err = v3.Repositories.CreateFile(c.ctx, owner, repo, path, opt)
	return err
}

func (c *Client) GetBlob(owner, repo, sha, path string) (Blob, error) {
	var q getFileSHAQuery
	err := c.g.Query(c.ctx, &q, map[string]interface{}{
		"owner":      githubv4.String(owner),
		"name":       githubv4.String(repo),
		"expression": githubv4.String(fmt.Sprintf("%s:%s", sha, path)),
	})
	if err != nil {
		return Blob{}, nil
	}
	return q.Repository.Object.Blob, nil
}
