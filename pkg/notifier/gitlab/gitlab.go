package gitlab

import (
	"fmt"

	gitlab "github.com/xanzy/go-gitlab"
)

// API is GitLab API interface
type API interface {
	CreateMergeRequestNote(mergeRequest int, opt *gitlab.CreateMergeRequestNoteOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Note, *gitlab.Response, error)
	UpdateMergeRequestNote(mergeRequest, note int, opt *gitlab.UpdateMergeRequestNoteOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Note, *gitlab.Response, error)
	ListMergeRequestNotes(mergeRequest int, opt *gitlab.ListMergeRequestNotesOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Note, *gitlab.Response, error)
	GetMergeRequest(mergeRequest int, opt *gitlab.GetMergeRequestsOptions, options ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error)
	UpdateMergeRequest(mergeRequest int, opt *gitlab.UpdateMergeRequestOptions, options ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error)
	PostCommitComment(sha string, opt *gitlab.PostCommitCommentOptions, options ...gitlab.RequestOptionFunc) (*gitlab.CommitComment, *gitlab.Response, error)
	ListCommits(opt *gitlab.ListCommitsOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Commit, *gitlab.Response, error)
	AddMergeRequestLabels(labels *[]string, mergeRequest int) (gitlab.Labels, error)
	RemoveMergeRequestLabels(labels *[]string, mergeRequest int) (gitlab.Labels, error)
	ListMergeRequestLabels(mergeRequest int, opt *gitlab.GetMergeRequestsOptions, options ...gitlab.RequestOptionFunc) (gitlab.Labels, error)
	GetLabel(labelName string, options ...gitlab.RequestOptionFunc) (*gitlab.Label, *gitlab.Response, error)
	UpdateLabel(opt *gitlab.UpdateLabelOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Label, *gitlab.Response, error)
	GetCommit(sha string, options ...gitlab.RequestOptionFunc) (*gitlab.Commit, *gitlab.Response, error)
}

// GitLab represents the attribute information necessary for requesting GitLab API
type GitLab struct {
	*gitlab.Client
	namespace, project string
}

// CreateMergeRequestNote is a wrapper of https://godoc.org/github.com/xanzy/go-gitlab#NotesService.CreateMergeRequestNote
func (g *GitLab) CreateMergeRequestNote(mergeRequest int, opt *gitlab.CreateMergeRequestNoteOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Note, *gitlab.Response, error) {
	return g.Client.Notes.CreateMergeRequestNote(fmt.Sprintf("%s/%s", g.namespace, g.project), mergeRequest, opt, options...)
}

// UpdateMergeRequestNote is a wrapper of https://pkg.go.dev/github.com/xanzy/go-gitlab#NotesService.UpdateMergeRequestNote
func (g *GitLab) UpdateMergeRequestNote(mergeRequest, note int, opt *gitlab.UpdateMergeRequestNoteOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Note, *gitlab.Response, error) {
	return g.Client.Notes.UpdateMergeRequestNote(fmt.Sprintf("%s/%s", g.namespace, g.project), mergeRequest, note, opt, options...)
}

// ListMergeRequestNotes is a wrapper of https://godoc.org/github.com/xanzy/go-gitlab#NotesService.ListMergeRequestNotes
func (g *GitLab) ListMergeRequestNotes(mergeRequest int, opt *gitlab.ListMergeRequestNotesOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Note, *gitlab.Response, error) {
	return g.Client.Notes.ListMergeRequestNotes(fmt.Sprintf("%s/%s", g.namespace, g.project), mergeRequest, opt, options...)
}

// GetMergerRequest is a wrapper of https://pkg.go.dev/github.com/xanzy/go-gitlab#MergeRequestsService.GetMergeRequest
func (g *GitLab) GetMergeRequest(mergeRequest int, opt *gitlab.GetMergeRequestsOptions, options ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error) {
	return g.Client.MergeRequests.GetMergeRequest(fmt.Sprintf("%s/%s", g.namespace, g.project), mergeRequest, opt, options...)
}

// UpdateMergerRequest is a wrapper of https://pkg.go.dev/github.com/xanzy/go-gitlab#MergeRequestsService.UpdateMergeRequest
func (g *GitLab) UpdateMergeRequest(mergeRequest int, opt *gitlab.UpdateMergeRequestOptions, options ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error) {
	return g.Client.MergeRequests.UpdateMergeRequest(fmt.Sprintf("%s/%s", g.namespace, g.project), mergeRequest, opt, options...)
}

// PostCommitComment is a wrapper of https://godoc.org/github.com/xanzy/go-gitlab#CommitsService.PostCommitComment
func (g *GitLab) PostCommitComment(sha string, opt *gitlab.PostCommitCommentOptions, options ...gitlab.RequestOptionFunc) (*gitlab.CommitComment, *gitlab.Response, error) {
	return g.Client.Commits.PostCommitComment(fmt.Sprintf("%s/%s", g.namespace, g.project), sha, opt, options...)
}

// ListCommits is a wrapper of https://godoc.org/github.com/xanzy/go-gitlab#CommitsService.ListCommits
func (g *GitLab) ListCommits(opt *gitlab.ListCommitsOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Commit, *gitlab.Response, error) {
	return g.Client.Commits.ListCommits(fmt.Sprintf("%s/%s", g.namespace, g.project), opt, options...)
}

// AddMergeRequestLabels adds labels on the merge request.
func (g *GitLab) AddMergeRequestLabels(labels *[]string, mergeRequest int) (gitlab.Labels, error) {
	var addLabels gitlab.Labels
	for _, label := range *labels {
		addLabels = append(addLabels, label)
	}

	updatedMergeRequest, _, err := g.Client.MergeRequests.UpdateMergeRequest(fmt.Sprintf("%s/%s", g.namespace, g.project), mergeRequest, &gitlab.UpdateMergeRequestOptions{AddLabels: &addLabels})
	if err != nil {
		return nil, err
	}
	return updatedMergeRequest.Labels, nil
}

// RemoveMergeRequestLabels removes labels on the merge request.
func (g *GitLab) RemoveMergeRequestLabels(labels *[]string, mergeRequest int) (gitlab.Labels, error) {
	var removeLabels gitlab.Labels
	for _, label := range *labels {
		removeLabels = append(removeLabels, label)
	}

	updatedMergeRequest, _, err := g.Client.MergeRequests.UpdateMergeRequest(fmt.Sprintf("%s/%s", g.namespace, g.project), mergeRequest, &gitlab.UpdateMergeRequestOptions{RemoveLabels: &removeLabels})
	if err != nil {
		return nil, err
	}
	return updatedMergeRequest.Labels, nil
}

// ListMergeRequestLabels lists labels on the merger request
func (g *GitLab) ListMergeRequestLabels(mergeRequest int, opt *gitlab.GetMergeRequestsOptions, options ...gitlab.RequestOptionFunc) (gitlab.Labels, error) {
	mr, _, err := g.Client.MergeRequests.GetMergeRequest(fmt.Sprintf("%s/%s", g.namespace, g.project), mergeRequest, opt, options...)
	if err != nil {
		return nil, err
	}
	return mr.Labels, nil
}

// GetLabel is a wrapper of https://pkg.go.dev/github.com/xanzy/go-gitlab#LabelsService.GetLabel
func (g *GitLab) GetLabel(labelName string, options ...gitlab.RequestOptionFunc) (*gitlab.Label, *gitlab.Response, error) {
	return g.Client.Labels.GetLabel(fmt.Sprintf("%s%s", g.namespace, g.project), labelName, options...)
}

// UpdateLabel is a wrapper of https://pkg.go.dev/github.com/xanzy/go-gitlab#LabelsService.UpdateLabel
func (g *GitLab) UpdateLabel(opt *gitlab.UpdateLabelOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Label, *gitlab.Response, error) {
	return g.Client.Labels.UpdateLabel(fmt.Sprintf("%s/%s", g.namespace, g.project), opt, options...)
}

// GetCommit is a wrapper of https://pkg.go.dev/github.com/xanzy/go-gitlab#CommitsService.GetCommit
func (g *GitLab) GetCommit(sha string, options ...gitlab.RequestOptionFunc) (*gitlab.Commit, *gitlab.Response, error) {
	return g.Client.Commits.GetCommit(fmt.Sprintf("%s/%s", g.namespace, g.project), sha, options...)
}
