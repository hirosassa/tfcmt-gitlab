package gitlab

import (
	"errors"
	"reflect"
	"testing"

	gitlabmock "github.com/hirosassa/tfcmt-gitlab/pkg/notifier/gitlab/gen"
	gitlab "github.com/xanzy/go-gitlab"
	"go.uber.org/mock/gomock"
)

func TestCommentPost(t *testing.T) {
	t.Parallel()
	body := "body"
	testCases := []struct {
		name                string
		config              Config
		createMockGitLabAPI func(ctrl *gomock.Controller) *gitlabmock.MockAPI
		body                string
		opt                 PostOptions
		ok                  bool
	}{
		{
			name:   "should post",
			config: newFakeConfig(),
			createMockGitLabAPI: func(ctrl *gomock.Controller) *gitlabmock.MockAPI {
				api := gitlabmock.NewMockAPI(ctrl)
				api.EXPECT().CreateMergeRequestNote(1, &gitlab.CreateMergeRequestNoteOptions{Body: gitlab.String(body)}).Return(&gitlab.Note{}, nil, nil)
				return api
			},
			body: body,
			opt: PostOptions{
				Number:   1,
				Revision: "abcd",
			},
			ok: true,
		},
		{
			name:   "should get mriid when PostOptions.Number is 0 and has PostOptions.Revision",
			config: newFakeConfig(),
			createMockGitLabAPI: func(ctrl *gomock.Controller) *gitlabmock.MockAPI {
				api := gitlabmock.NewMockAPI(ctrl)
				mriid := 1
				api.EXPECT().ListMergeRequestsByCommit("abcd").Return([]*gitlab.MergeRequest{{IID: mriid}}, nil, nil)
				api.EXPECT().CreateMergeRequestNote(mriid, &gitlab.CreateMergeRequestNoteOptions{Body: gitlab.String(body)}).Return(&gitlab.Note{}, nil, nil)
				return api
			},
			body: body,
			opt: PostOptions{
				Number:   0,
				Revision: "abcd",
			},
			ok: true,
		},
		{
			name:   "should post number 2",
			config: newFakeConfig(),
			createMockGitLabAPI: func(ctrl *gomock.Controller) *gitlabmock.MockAPI {
				api := gitlabmock.NewMockAPI(ctrl)
				api.EXPECT().CreateMergeRequestNote(2, &gitlab.CreateMergeRequestNoteOptions{Body: gitlab.String(body)}).Return(&gitlab.Note{}, nil, nil)
				return api
			},
			body: body,
			opt: PostOptions{
				Number:   2,
				Revision: "",
			},
			ok: true,
		},
		{
			name:   "should error PostOptions number=0 and Revision is empty",
			config: newFakeConfig(),
			createMockGitLabAPI: func(ctrl *gomock.Controller) *gitlabmock.MockAPI {
				api := gitlabmock.NewMockAPI(ctrl)
				return api
			},
			body: "",
			opt: PostOptions{
				Number:   0,
				Revision: "",
			},
			ok: false,
		},
		{
			name:   "should postForRevision when listMergeRequestIIDs is failed",
			config: newFakeConfig(),
			createMockGitLabAPI: func(ctrl *gomock.Controller) *gitlabmock.MockAPI {
				api := gitlabmock.NewMockAPI(ctrl)
				api.EXPECT().ListMergeRequestsByCommit("revision").Return(nil, nil, errors.New("error"))
				// PostCommitComment should be called
				api.EXPECT().PostCommitComment("revision", &gitlab.PostCommitCommentOptions{Note: gitlab.String(body)}).Return(&gitlab.CommitComment{}, nil, nil)
				return api
			},
			body: body,
			opt: PostOptions{
				Number:   0,
				Revision: "revision",
			},
			ok: true,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			client, err := NewClient(testCase.config)
			if err != nil {
				t.Fatal(err)
			}

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			client.API = testCase.createMockGitLabAPI(mockCtrl)

			err = client.Comment.Post(testCase.body, testCase.opt)
			if (err == nil) != testCase.ok {
				t.Errorf("got error %q", err)
			}
		})
	}
}

func TestCommentList(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name                string
		config              Config
		createMockGitLabAPI func(ctrl *gomock.Controller) *gitlabmock.MockAPI
		number              int
		ok                  bool
		comments            []*gitlab.Note
	}{
		{
			name:   "should list comments",
			config: newFakeConfig(),
			createMockGitLabAPI: func(ctrl *gomock.Controller) *gitlabmock.MockAPI {
				api := gitlabmock.NewMockAPI(ctrl)
				api.EXPECT().ListMergeRequestNotes(1, gomock.Any()).MaxTimes(2).DoAndReturn(
					func(mergeRequest int, opt *gitlab.ListMergeRequestNotesOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Note, *gitlab.Response, error) {
						res := []*gitlab.Note{
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

						return res, resp, nil
					},
				)

				return api
			},
			number: 1,
			ok:     true,
			comments: []*gitlab.Note{
				// page1
				{
					ID:   371748792,
					Body: "comment 1",
				},
				{
					ID:   371765743,
					Body: "comment 2",
				},
				// page2
				{
					ID:   371748792,
					Body: "comment 1",
				},
				{
					ID:   371765743,
					Body: "comment 2",
				},
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			client, err := NewClient(testCase.config)
			if err != nil {
				t.Fatal(err)
			}

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			client.API = testCase.createMockGitLabAPI(mockCtrl)

			comments, err := client.Comment.List(testCase.number)
			if (err == nil) != testCase.ok {
				t.Errorf("got error %q", err)
			}
			if !reflect.DeepEqual(comments, testCase.comments) {
				t.Errorf("got %v but want %v", comments, testCase.comments)
			}
		})
	}
}

func TestCommentListSentinel(t *testing.T) {
	t.Parallel()

	client, err := NewClient(newFakeConfig())
	if err != nil {
		t.Fatal(err)
	}

	createMockGitLabAPI := func(ctrl *gomock.Controller) *gitlabmock.MockAPI {
		api := gitlabmock.NewMockAPI(ctrl)
		api.EXPECT().ListMergeRequestNotes(1, gomock.Any()).MaxTimes(100).Return(
			[]*gitlab.Note{},
			&gitlab.Response{
				NextPage: 1, // just cause infinite loop
			},
			nil,
		)

		return api
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	client.API = createMockGitLabAPI(mockCtrl)

	_, err = client.Comment.List(1) // no assert res, only assert `.MaxTimes(100)`
	if err != nil {
		t.Errorf("got error %q", err)
	}
}
