package gitlab

import (
	"testing"

	"github.com/hirosassa/tfcmt-gitlab/pkg/notifier"
	"github.com/hirosassa/tfcmt-gitlab/pkg/terraform"
)

func TestNotifyNotify(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name      string
		config    Config
		ok        bool
		exitCode  int
		paramExec notifier.ParamExec
	}{
		{
			name: "invalid body (cannot parse)",
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
		client, err := NewClient(testCase.config)
		if err != nil {
			t.Fatal(err)
		}
		api := newFakeAPI()
		client.API = &api
		exitCode, err := client.Notify.Notify(testCase.paramExec)
		if (err == nil) != testCase.ok {
			t.Errorf("test case: %s, got error %q", testCase.name, err)
		}
		if exitCode != testCase.exitCode {
			t.Errorf("test case: %s, got %q but want %q", testCase.name, exitCode, testCase.exitCode)
		}
	}
}
