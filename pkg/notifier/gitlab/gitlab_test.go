package gitlab

import (
	"github.com/hirosassa/tfcmt-gitlab/pkg/terraform"
)

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
