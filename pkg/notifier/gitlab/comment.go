package gitlab

import (
	"fmt"

	"github.com/sirupsen/logrus"
	gitlab "github.com/xanzy/go-gitlab"
)

const (
	listPerPage = 100
	maxPages    = 100
)

// CommentService handles communication with the comment related
// methods of GitLab API
type CommentService service

// PostOptions specifies the optional parameters to post comments to a pull request
type PostOptions struct {
	Number   int
	Revision string
}

// Post posts comment
func (g *CommentService) Post(body string, opt PostOptions) error {
	if opt.Number == 0 && opt.Revision == "" {
		return fmt.Errorf("gitlab.comment.post: Number or Revision is required")
	}

	if opt.Number == 0 {
		mrs, err := g.client.Commits.ListMergeRequestIIDsByRevision(opt.Revision)
		if err != nil || len(mrs) == 0 {
			return g.postForRevision(body, opt.Revision)
		}

		// Rewrite the MR number to the first MR which is associated with revision.
		opt.Number = mrs[0]
	}

	_, _, err := g.client.API.CreateMergeRequestNote(
		opt.Number,
		&gitlab.CreateMergeRequestNoteOptions{Body: gitlab.String(body)},
	)
	return err
}

func (g *CommentService) postForRevision(body, revision string) error {
	_, _, err := g.client.API.PostCommitComment(
		revision,
		&gitlab.PostCommitCommentOptions{Note: gitlab.String(body)},
	)
	return err
}

// Patch patches on the specific comment
func (g *CommentService) Patch(note int, body string, opt PostOptions) error {
	_, _, err := g.client.API.UpdateMergeRequestNote(
		opt.Number,
		note,
		&gitlab.UpdateMergeRequestNoteOptions{Body: gitlab.String(body)},
	)
	return err
}

// List lists comments on GitLab merge requests
func (g *CommentService) List(number int) ([]*gitlab.Note, error) {
	allComments := make([]*gitlab.Note, 0)

	opt := &gitlab.ListMergeRequestNotesOptions{
		ListOptions: gitlab.ListOptions{
			Page:    1,
			PerPage: listPerPage,
		},
	}

	for sentinel := 1; ; sentinel++ {
		comments, resp, err := g.client.API.ListMergeRequestNotes(
			number,
			opt,
		)
		if err != nil {
			return nil, err
		}

		allComments = append(allComments, comments...)

		if resp.NextPage == 0 {
			break
		}

		if sentinel > maxPages {
			logE := logrus.WithFields(logrus.Fields{
				"program": "tfcmt",
			})
			logE.WithField("maxPages", maxPages).Debug("gitlab.comment.list: too many pages, something went wrong")
			break
		}

		opt.Page = resp.NextPage
	}

	return allComments, nil
}
