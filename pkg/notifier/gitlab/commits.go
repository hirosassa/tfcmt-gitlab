package gitlab

import (
	"errors"
)

// CommitsService handles communication with the commits related
// methods of GitLab API
type CommitsService service

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
		result[i] = int(mr.IID)
	}
	return result, nil
}
