package gitlab

import (
	"testing"

	"github.com/hirosassa/tfcmt-gitlab/pkg/notifier"
	gitlabmock "github.com/hirosassa/tfcmt-gitlab/pkg/notifier/gitlab/gen"
	"github.com/hirosassa/tfcmt-gitlab/pkg/terraform"
	gitlab "github.com/xanzy/go-gitlab"
	"go.uber.org/mock/gomock"
)

func TestNotifyNotify(t *testing.T) { //nolint:maintidx
	t.Parallel()
	testCases := []struct {
		name                string
		createMockGitLabAPI func(ctrl *gomock.Controller) *gitlabmock.MockAPI
		config              Config
		ok                  bool
		exitCode            int
		paramExec           notifier.ParamExec
	}{
		{
			name: "invalid body (cannot parse)",
			createMockGitLabAPI: func(ctrl *gomock.Controller) *gitlabmock.MockAPI {
				api := gitlabmock.NewMockAPI(ctrl)
				api.EXPECT().CreateMergeRequestNote(1, gomock.Any()).Return(nil, nil, nil)
				return api
			},
			config: Config{
				Token:     "token",
				NameSpace: "namespace",
				Project:   "project",
				MR: MergeRequest{
					Revision: "abcd",
					Number:   1,
					Message:  "message",
				},
				Parser:             terraform.NewPlanParser(),
				Template:           terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
			},
			paramExec: notifier.ParamExec{
				Stdout:   "body",
				ExitCode: 1,
			},
			ok:       true,
			exitCode: 1,
		},
		{
			name: "invalid mr",
			createMockGitLabAPI: func(ctrl *gomock.Controller) *gitlabmock.MockAPI {
				api := gitlabmock.NewMockAPI(ctrl)
				return api
			},
			config: Config{
				Token:     "token",
				NameSpace: "namespace",
				Project:   "project",
				MR: MergeRequest{
					Revision: "",
					Number:   0,
					Message:  "message",
				},
				Parser:             terraform.NewPlanParser(),
				Template:           terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
			},
			paramExec: notifier.ParamExec{
				Stdout:   "Plan 1 to add",
				ExitCode: 0,
			},
			ok:       false,
			exitCode: 0,
		},
		{
			name: "valid, error",
			createMockGitLabAPI: func(ctrl *gomock.Controller) *gitlabmock.MockAPI {
				api := gitlabmock.NewMockAPI(ctrl)
				api.EXPECT().CreateMergeRequestNote(1, gomock.Any()).Return(nil, nil, nil)
				return api
			},
			config: Config{
				Token:     "token",
				NameSpace: "namespace",
				Project:   "project",
				MR: MergeRequest{
					Revision: "",
					Number:   1,
					Message:  "message",
				},
				Parser:             terraform.NewPlanParser(),
				Template:           terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
			},
			paramExec: notifier.ParamExec{
				Stdout:   "Error: hoge",
				ExitCode: 1,
			},
			ok:       true,
			exitCode: 1,
		},
		{
			name: "valid, and isMR",
			createMockGitLabAPI: func(ctrl *gomock.Controller) *gitlabmock.MockAPI {
				api := gitlabmock.NewMockAPI(ctrl)
				api.EXPECT().CreateMergeRequestNote(1, gomock.Any()).Return(nil, nil, nil)
				return api
			},
			config: Config{
				Token:     "token",
				NameSpace: "namespace",
				Project:   "project",
				MR: MergeRequest{
					Revision: "",
					Number:   1,
					Message:  "message",
				},
				Parser:             terraform.NewPlanParser(),
				Template:           terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
			},
			paramExec: notifier.ParamExec{
				Stdout:   "Plan 1 to add",
				ExitCode: 2,
			},
			ok:       true,
			exitCode: 2,
		},
		{
			name: "valid, and isRevision",
			createMockGitLabAPI: func(ctrl *gomock.Controller) *gitlabmock.MockAPI {
				api := gitlabmock.NewMockAPI(ctrl)
				api.EXPECT().ListMergeRequestsByCommit("revision-revision").Return([]*gitlab.MergeRequest{{IID: 1}}, nil, nil)
				api.EXPECT().CreateMergeRequestNote(1, gomock.Any()).Return(nil, nil, nil)
				return api
			},
			config: Config{
				Token:     "token",
				NameSpace: "namespace",
				Project:   "project",
				MR: MergeRequest{
					Revision: "revision-revision",
					Number:   0,
					Message:  "message",
				},
				Parser:             terraform.NewPlanParser(),
				Template:           terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
			},
			paramExec: notifier.ParamExec{
				Stdout:   "Plan 1 to add",
				ExitCode: 2,
			},
			ok:       true,
			exitCode: 2,
		},
		{
			name: "valid, and contains destroy",
			createMockGitLabAPI: func(ctrl *gomock.Controller) *gitlabmock.MockAPI {
				api := gitlabmock.NewMockAPI(ctrl)
				api.EXPECT().CreateMergeRequestNote(1, gomock.Any()).Return(nil, nil, nil)
				return api
			},
			config: Config{
				Token:     "token",
				NameSpace: "namespace",
				Project:   "project",
				MR: MergeRequest{
					Revision: "",
					Number:   1,
				},
				Parser:             terraform.NewPlanParser(),
				Template:           terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
			},
			paramExec: notifier.ParamExec{
				Stdout:   "Plan: 1 to add, 1 to destroy",
				ExitCode: 2,
			},
			ok:       true,
			exitCode: 2,
		},
		{
			name: "valid with no change",
			createMockGitLabAPI: func(ctrl *gomock.Controller) *gitlabmock.MockAPI {
				api := gitlabmock.NewMockAPI(ctrl)
				api.EXPECT().CreateMergeRequestNote(1, gomock.Any()).Return(nil, nil, nil)
				return api
			},
			config: Config{
				Token:     "token",
				NameSpace: "namespace",
				Project:   "project",
				MR: MergeRequest{
					Revision: "",
					Number:   1,
				},
				Parser:             terraform.NewPlanParser(),
				Template:           terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
			},
			paramExec: notifier.ParamExec{
				Stdout:   "No changes. Infrastructure is up-to-date.",
				ExitCode: 0,
			},
			ok:       true,
			exitCode: 0,
		},
		{
			name: "valid, contains destroy, but not to notify",
			createMockGitLabAPI: func(ctrl *gomock.Controller) *gitlabmock.MockAPI {
				api := gitlabmock.NewMockAPI(ctrl)
				api.EXPECT().CreateMergeRequestNote(1, gomock.Any()).Return(nil, nil, nil)
				return api
			},
			config: Config{
				Token:     "token",
				NameSpace: "namespace",
				Project:   "project",
				MR: MergeRequest{
					Revision: "",
					Number:   1,
				},
				Parser:             terraform.NewPlanParser(),
				Template:           terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
			},
			paramExec: notifier.ParamExec{
				Stdout:   "Plan: 1 to add, 1 to destroy",
				ExitCode: 2,
			},
			ok:       true,
			exitCode: 2,
		},
		{
			name: "apply case",
			createMockGitLabAPI: func(ctrl *gomock.Controller) *gitlabmock.MockAPI {
				api := gitlabmock.NewMockAPI(ctrl)
				api.EXPECT().ListCommits(gomock.Any()).Return([]*gitlab.Commit{{ID: "1"}, {ID: "2"}}, nil, nil)
				api.EXPECT().ListMergeRequestsByCommit("2").Return([]*gitlab.MergeRequest{{IID: 1}}, nil, nil)
				api.EXPECT().CreateMergeRequestNote(1, gomock.Any()).Return(nil, nil, nil)
				return api
			},
			config: Config{
				Token:     "token",
				NameSpace: "namespace",
				Project:   "project",
				MR: MergeRequest{
					Revision: "revision",
					Number:   0, // For apply, it is always 0
					Message:  "message",
				},
				Parser:             terraform.NewApplyParser(),
				Template:           terraform.NewApplyTemplate(terraform.DefaultApplyTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
			},
			paramExec: notifier.ParamExec{
				Stdout:   "Apply complete!",
				ExitCode: 0,
			},
			ok:       true,
			exitCode: 0,
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

			exitCode, err := client.Notify.Notify(testCase.paramExec)
			if (err == nil) != testCase.ok {
				t.Errorf("test case: %s, got error %q", testCase.name, err)
			}
			if exitCode != testCase.exitCode {
				t.Errorf("test case: %s, got %q but want %q", testCase.name, exitCode, testCase.exitCode)
			}
		})
	}
}
