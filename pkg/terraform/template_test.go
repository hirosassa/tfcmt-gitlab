package terraform_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hirosassa/tfcmt-gitlab/pkg/terraform"
)

func TestTemplate_Execute(t *testing.T) {
	t.Parallel()
	templ := terraform.NewPlanTemplate("")

	templ.SetValue(terraform.CommonTemplate{
		ExitCode: 0,
		Vars: map[string]string{
			"target": "test",
		},
		Link: "http://example.com",
	})

	got, err := templ.Execute()
	if err != nil {
		t.Fatal(err)
	}

	expect := `
## Plan Result (test)

[CI link](http://example.com)








`
	if diff := cmp.Diff(expect, got); diff != "" {
		t.Errorf("Template.Execute result diff (-expect, +got)\n%s", diff)
	}
}

func TestTemplate_IsSamePlan(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		oldStr         string
		commonTemplate terraform.CommonTemplate
		expect         bool
	}{
		{
			name: "should return true when plan succeeded",
			oldStr: `
## Plan Result (test)

[CI link](http://example.com)
`,
			commonTemplate: terraform.CommonTemplate{
				ExitCode: 0,
				Vars: map[string]string{
					"target": "test",
				},
			},
			expect: true,
		},
		{
			name: "should return true when plan failed",
			oldStr: `
## Plan Result (test)

[CI link](http://example.com)
`,
			commonTemplate: terraform.CommonTemplate{
				ExitCode: 1,
				Vars: map[string]string{
					"target": "test",
				},
			},
			expect: true,
		},
		{
			name: "should return false when oldStr is Apply Result",
			oldStr: `
## Apply Result (test)

[CI link](http://example.com)
`,
			commonTemplate: terraform.CommonTemplate{
				ExitCode: 0,
				Vars: map[string]string{
					"target": "test",
				},
			},
			expect: false,
		},
		{
			name: "should return false when target is empty",
			oldStr: `
## Plan Result

[CI link](http://example.com)
`,
			commonTemplate: terraform.CommonTemplate{
				ExitCode: 0,
				Vars:     nil,
			},
			expect: false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			templ := terraform.NewPlanTemplate("")
			templ.SetValue(testCase.commonTemplate)

			got := templ.IsSamePlan(testCase.oldStr)
			if testCase.expect != got {
				t.Errorf("Template.IsSamePlan should return %t, but got %t", testCase.expect, got)
			}
		})
	}
}
