package client

import (
	"bytes"
	"encoding/base64"

	"github.com/google/go-github/v39/github"
)

func (c *Client) CreateOrUpdate(owner, repo, sha, path, message string, content []byte) error {
	b, err := c.getBlob(owner, repo, sha, path)
	if err != nil {
		return err
	}

	if b != nil {
		preContent, err := base64.StdEncoding.DecodeString(*b.Content)
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
	if b != nil {
		opt.SHA = b.SHA
	}
	_, _, err = c.g.Repositories.CreateFile(c.ctx, owner, repo, path, opt)
	return err
}

func (c *Client) getBlob(owner, repo, sha, path string) (*github.Blob, error) {
	tree, _, err := c.g.Git.GetTree(c.ctx, owner, repo, sha, true)
	if err != nil {
		return nil, err
	}
	for _, ent := range tree.Entries {
		if *ent.Path == path {
			blob, _, err := c.g.Git.GetBlob(c.ctx, owner, repo, ent.GetSHA())
			if err != nil {
				return nil, err
			}
			return blob, nil
		}
	}
	return nil, nil
}
