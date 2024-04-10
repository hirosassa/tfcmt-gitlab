package gitlab

import (
	"github.com/hirosassa/tfcmt-gitlab/pkg/terraform"
)

type fakeAPI struct {
	API
	FakeCreateMergeRequestNote    func(mergeRequest int, opt *gitlab.CreateMergeRequestNoteOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Note, *gitlab.Response, error)
	FakeListMergeRequestNotes     func(mergeRequest int, opt *gitlab.ListMergeRequestNotesOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Note, *gitlab.Response, error)
	FakePostCommitComment         func(sha string, opt *gitlab.PostCommitCommentOptions, options ...gitlab.RequestOptionFunc) (*gitlab.CommitComment, *gitlab.Response, error)
	FakeListCommits               func(opt *gitlab.ListCommitsOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Commit, *gitlab.Response, error)
	FakeListMergeRequestsByCommit func(sha string, options ...gitlab.RequestOptionFunc) ([]*gitlab.MergeRequest, *gitlab.Response, error)
}

func (g *fakeAPI) CreateMergeRequestNote(mergeRequest int, opt *gitlab.CreateMergeRequestNoteOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Note, *gitlab.Response, error) {
	return g.FakeCreateMergeRequestNote(mergeRequest, opt, options...)
}

func (g *fakeAPI) ListMergeRequestNotes(mergeRequest int, opt *gitlab.ListMergeRequestNotesOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Note, *gitlab.Response, error) {
	return g.FakeListMergeRequestNotes(mergeRequest, opt, options...)
}

func (g *fakeAPI) PostCommitComment(sha string, opt *gitlab.PostCommitCommentOptions, options ...gitlab.RequestOptionFunc) (*gitlab.CommitComment, *gitlab.Response, error) {
	return g.FakePostCommitComment(sha, opt, options...)
}

func (g *fakeAPI) ListCommits(opt *gitlab.ListCommitsOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Commit, *gitlab.Response, error) {
	return g.FakeListCommits(opt, options...)
}

func (g *fakeAPI) ListMergeRequestsByCommit(sha string, options ...gitlab.RequestOptionFunc) ([]*gitlab.MergeRequest, *gitlab.Response, error) {
	return g.FakeListMergeRequestsByCommit(sha, options...)
}

func newFakeAPI() fakeAPI {
	return fakeAPI{
		FakeCreateMergeRequestNote: func(mergeRequest int, opt *gitlab.CreateMergeRequestNoteOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Note, *gitlab.Response, error) {
			return &gitlab.Note{
				ID:   371748792,
				Body: "comment 1",
			}, nil, nil
		},
		FakeListMergeRequestNotes: func(mergeRequest int, opt *gitlab.ListMergeRequestNotesOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Note, *gitlab.Response, error) {
			var comments []*gitlab.Note
			comments = []*gitlab.Note{
				// same response to any page for now
				{
					ID:   371748792,
					Body: "comment 1",
				},
				{
					ID:   371765743,
					Body: "comment 2",
				},
			}

			// fake pagination with 2 pages
			resp := &gitlab.Response{
				NextPage: 0,
			}
			if opt.Page == 1 {
				resp.NextPage = 2
			}

			return comments, resp, nil
		},
		FakePostCommitComment: func(sha string, opt *gitlab.PostCommitCommentOptions, options ...gitlab.RequestOptionFunc) (*gitlab.CommitComment, *gitlab.Response, error) {
			return &gitlab.CommitComment{
				Note: "comment 1",
			}, nil, nil
		},
		FakeListCommits: func(opt *gitlab.ListCommitsOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Commit, *gitlab.Response, error) {
			var commits []*gitlab.Commit
			commits = []*gitlab.Commit{
				{
					ID: "04e0917e448b662c2b16330fad50e97af16ff27a",
				},
				{
					ID: "04e0917e448b662c2b16330fad50e97af16ff27b",
				},
				{
					ID: "04e0917e448b662c2b16330fad50e97af16ff27c",
				},
			}
			return commits, nil, nil
		},
		FakeListMergeRequestsByCommit: func(sha string, options ...gitlab.RequestOptionFunc) ([]*gitlab.MergeRequest, *gitlab.Response, error) {
			return []*gitlab.MergeRequest{
				{
					IID: 1,
				},
			}, nil, nil
		},
	}
}

func newFakeConfig() Config {
	return Config{
		Token:     "token",
		NameSpace: "owner",
		Project:   "repo",
		MR: MergeRequest{
			Revision: "abcd",
			Number:   1,
			Message:  "message",
		},
		Parser:   terraform.NewPlanParser(),
		Template: terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
	}
}
