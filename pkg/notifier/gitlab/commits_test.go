package gitlab

import (
	"testing"

	gitlabmock "github.com/hirosassa/tfcmt-gitlab/pkg/notifier/gitlab/gen"
	gitlab "github.com/xanzy/go-gitlab"
	"go.uber.org/mock/gomock"
)

func TestCommitsList(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name                string
		createMockGitLabAPI func(ctrl *gomock.Controller) *gitlabmock.MockAPI
		revision            string
		ok                  bool
	}{
		{
			name: "should list commits",
			createMockGitLabAPI: func(ctrl *gomock.Controller) *gitlabmock.MockAPI {
				api := gitlabmock.NewMockAPI(ctrl)
				api.EXPECT().ListCommits(gomock.Cond(func(x any) bool {
					_, ok := x.(*gitlab.ListCommitsOptions)
					return ok
				})).Return([]*gitlab.Commit{}, nil, nil)
				return api
			},
			revision: "04e0917e448b662c2b16330fad50e97af16ff27a",
			ok:       true,
		},
		{
			name: "should return error when revision is empty",
			createMockGitLabAPI: func(ctrl *gomock.Controller) *gitlabmock.MockAPI {
				api := gitlabmock.NewMockAPI(ctrl)
				return api
			},
			revision: "",
			ok:       false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			cfg := newFakeConfig()
			client, err := NewClient(cfg)
			if err != nil {
				t.Fatal(err)
			}

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			client.API = testCase.createMockGitLabAPI(mockCtrl)

			_, err = client.Commits.List(testCase.revision)
			if (err == nil) != testCase.ok {
				t.Errorf("got error %q", err)
			}
		})
	}
}

func TestCommitsLastOne(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		commits  []string
		revision string
		lastRev  string
		ok       bool
	}{
		{
			// ok
			commits: []string{
				"04e0917e448b662c2b16330fad50e97af16ff27a",
				"04e0917e448b662c2b16330fad50e97af16ff27b",
				"04e0917e448b662c2b16330fad50e97af16ff27c",
			},
			revision: "04e0917e448b662c2b16330fad50e97af16ff27a",
			lastRev:  "04e0917e448b662c2b16330fad50e97af16ff27b",
			ok:       true,
		},
		{
			// no revision
			commits: []string{
				"04e0917e448b662c2b16330fad50e97af16ff27a",
				"04e0917e448b662c2b16330fad50e97af16ff27b",
				"04e0917e448b662c2b16330fad50e97af16ff27c",
			},
			revision: "",
			lastRev:  "",
			ok:       false,
		},
		{
			// no commits
			commits:  []string{},
			revision: "04e0917e448b662c2b16330fad50e97af16ff27a",
			lastRev:  "",
			ok:       false,
		},
	}

	for _, testCase := range testCases {
		cfg := newFakeConfig()
		client, err := NewClient(cfg)
		if err != nil {
			t.Fatal(err)
		}
		commit, err := client.Commits.lastOne(testCase.commits, testCase.revision)
		if (err == nil) != testCase.ok {
			t.Errorf("got error %q", err)
		}
		if commit != testCase.lastRev {
			t.Errorf("got %q but want %q", commit, testCase.lastRev)
		}
	}
}
