package gitlab

import (
	"errors"

	gitlab "github.com/xanzy/go-gitlab"
)

// CommitsService handles communication with the commits related
// methods of GitLab API
type CommitsService service

// List lists commits on a repository
func (g *CommitsService) List(revision string) ([]string, error) {
	if revision == "" {
		return []string{}, errors.New("no revision specified")
	}
	commits, _, err := g.client.API.ListCommits(
		&gitlab.ListCommitsOptions{},
	)
	if err != nil {
		return nil, err
	}
	s := make([]string, len(commits))
	for i, commit := range commits {
		s[i] = commit.ID
	}
	return s, nil
}

func (g *CommitsService) ListMergeRequestIIDsByRevision(revision string) ([]int, error) {
	if revision == "" {
		return nil, errors.New("no revision specified")
	}
	mrs, _, err := g.client.API.ListMergeRequestsByCommit(revision)
	if err != nil {
		return nil, err
	}

	result := make([]int, len(mrs))
	for i, mr := range mrs {
		result[i] = mr.IID
	}
	return result, nil
}

// lastOne returns the hash of the previous commit of the given commit
func (g *CommitsService) lastOne(commits []string, revision string) (string, error) {
	if revision == "" {
		return "", errors.New("no revision specified")
	}
	if len(commits) == 0 {
		return "", errors.New("no commits")
	}

	return commits[1], nil
}
