package gitlab

import (
	"fmt"

	gitlab "github.com/xanzy/go-gitlab"
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
	comments, _, err := g.client.API.ListMergeRequestNotes(
		number,
		&gitlab.ListMergeRequestNotesOptions{},
	)
	return comments, err
}
