package terraform

import (
	"bytes"
	htmltemplate "html/template"
	"strings"
	texttemplate "text/template"

	"github.com/Masterminds/sprig/v3"
)

const (
	// DefaultPlanTemplate is a default template for terraform plan
	DefaultPlanTemplate = `
{{template "plan_title" .}}

{{if .Link}}[CI link]({{.Link}}){{end}}

{{template "deletion_warning" .}}
{{template "result" .}}
{{template "updated_resources" .}}

{{template "changed_result" .}}
{{template "change_outside_terraform" .}}
{{template "warning" .}}
{{template "error_messages" .}}`

	// DefaultApplyTemplate is a default template for terraform apply
	DefaultApplyTemplate = `
{{template "apply_title" .}}

{{if .Link}}[CI link]({{.Link}}){{end}}

{{if ne .ExitCode 0}}{{template "guide_apply_failure" .}}{{end}}

{{template "result" .}}

<details><summary>Details (Click me)</summary>
{{wrapCode .CombinedOutput}}
</details>
{{template "error_messages" .}}`

	// DefaultPlanParseErrorTemplate is a default template for terraform plan parse error
	DefaultPlanParseErrorTemplate = `
{{template "plan_title" .}}

{{if .Link}}[CI link]({{.Link}}){{end}}

It failed to parse the result.

<details><summary>Details (Click me)</summary>
{{wrapCode .CombinedOutput}}
</details>
`

	// DefaultApplyParseErrorTemplate  is a default template for terraform apply parse error
	DefaultApplyParseErrorTemplate = `
{{template "apply_title" .}}

{{if .Link}}[CI link]({{.Link}}){{end}}

{{template "guide_apply_parse_error" .}}

It failed to parse the result.

<details><summary>Details (Click me)</summary>
{{wrapCode .CombinedOutput}}
</details>
`

	planTitleTemplate = "## {{if eq .ExitCode 1}}:x: Plan Failed{{else}}Plan Result{{end}}{{if .Vars.target}} ({{.Vars.target}}){{end}}"

	applyTitleTemplate = "## {{if eq .ExitCode 0}}:white_check_mark: Apply Succeeded{{else}}:x: Apply Failed{{end}}{{if .Vars.target}} ({{.Vars.target}}){{end}}"

	resultTemplate = "{{if .Result}}<pre><code>{{ .Result }}</code></pre>{{end}}"

	updatedResourcesTemplate = `{{if .CreatedResources}}
* Create
{{- range .CreatedResources}}
  * {{.}}
{{- end}}{{end}}{{if .UpdatedResources}}
* Update
{{- range .UpdatedResources}}
  * {{.}}
{{- end}}{{end}}{{if .DeletedResources}}
* Delete
{{- range .DeletedResources}}
  * {{.}}
{{- end}}{{end}}{{if .ReplacedResources}}
* Replace
{{- range .ReplacedResources}}
  * {{.}}
{{- end}}{{end}}`

	deletionWarningTemplate = `{{if .HasDestroy}}
### :warning: Resource Deletion will happen :warning:
This plan contains resource delete operation. Please check the plan result very carefully!
{{end}}`

	changedResultTemplate = `{{if .ChangedResult}}
<details><summary>Change Result (Click me)</summary>
{{wrapCode .ChangedResult}}
</details>
{{end}}`

	changeOutsideTerraformTemplate = `{{if .ChangeOutsideTerraform}}
<details><summary>:information_source: Objects have changed outside of Terraform</summary>

_This feature was introduced from [Terraform v0.15.4](https://github.com/hashicorp/terraform/releases/tag/v0.15.4)._
{{wrapCode .ChangeOutsideTerraform}}
</details>
{{end}}`

	warningTemplate = `{{if .Warning}}
## :warning: Warnings :warning:
{{wrapCode .Warning}}
{{end}}`

	errorMessagesTemplate = `{{if .ErrorMessages}}
## :warning: Errors
{{range .ErrorMessages}}
* {{. -}}
{{- end}}{{end}}`
)

// CommonTemplate represents template entities
type CommonTemplate struct {
	Result                 string
	ChangedResult          string
	ChangeOutsideTerraform string
	Warning                string
	Link                   string
	UseRawOutput           bool
	HasDestroy             bool
	Vars                   map[string]string
	Templates              map[string]string
	Stdout                 string
	Stderr                 string
	CombinedOutput         string
	ExitCode               int
	ErrorMessages          []string
	CreatedResources       []string
	UpdatedResources       []string
	DeletedResources       []string
	ReplacedResources      []string
}

// Template is a default template for terraform commands
type Template struct {
	Template string
	CommonTemplate
}

// NewPlanTemplate is PlanTemplate initializer
func NewPlanTemplate(template string) *Template {
	if template == "" {
		template = DefaultPlanTemplate
	}
	return &Template{
		Template: template,
	}
}

// NewApplyTemplate is ApplyTemplate initializer
func NewApplyTemplate(template string) *Template {
	if template == "" {
		template = DefaultApplyTemplate
	}
	return &Template{
		Template: template,
	}
}

func NewPlanParseErrorTemplate(template string) *Template {
	if template == "" {
		template = DefaultPlanParseErrorTemplate
	}
	return &Template{
		Template: template,
	}
}

func NewApplyParseErrorTemplate(template string) *Template {
	if template == "" {
		template = DefaultApplyParseErrorTemplate
	}
	return &Template{
		Template: template,
	}
}

func avoidHTMLEscape(text string) htmltemplate.HTML {
	return htmltemplate.HTML(text) //nolint:gosec
}

func wrapCode(text string) interface{} {
	if len(text) > 60000 { //nolint:gomnd
		text = text[:20000] + `

# ...
# ... The maximum length of GitHub Comment is 65536, so the content is omitted by tfcmt.
# ...

` + text[len(text)-20000:]
	}
	if strings.Contains(text, "```") {
		return `<pre><code>` + text + `</code></pre>`
	}
	return htmltemplate.HTML("\n```hcl\n" + text + "\n```\n") //nolint:gosec
}

func generateOutput(kind, template string, data any, useRawOutput bool) (string, error) {
	var b bytes.Buffer

	if useRawOutput {
		tpl, err := texttemplate.New(kind).Funcs(texttemplate.FuncMap{
			"avoidHTMLEscape": avoidHTMLEscape,
			"wrapCode":        wrapCode,
		}).Funcs(sprig.TxtFuncMap()).Parse(template)
		if err != nil {
			return "", err
		}
		if err := tpl.Execute(&b, data); err != nil {
			return "", err
		}
	} else {
		tpl, err := htmltemplate.New(kind).Funcs(htmltemplate.FuncMap{
			"avoidHTMLEscape": avoidHTMLEscape,
			"wrapCode":        wrapCode,
		}).Funcs(sprig.FuncMap()).Parse(template)
		if err != nil {
			return "", err
		}
		if err := tpl.Execute(&b, data); err != nil {
			return "", err
		}
	}

	return b.String(), nil
}

// Execute binds the execution result of terraform command into template
func (t *Template) Execute() (string, error) {
	templates := map[string]string{
		"plan_title":               planTitleTemplate,
		"apply_title":              applyTitleTemplate,
		"result":                   resultTemplate,
		"updated_resources":        updatedResourcesTemplate,
		"deletion_warning":         deletionWarningTemplate,
		"changed_result":           changedResultTemplate,
		"change_outside_terraform": changeOutsideTerraformTemplate,
		"warning":                  warningTemplate,
		"error_messages":           errorMessagesTemplate,
		"guide_apply_failure":      "",
		"guide_apply_parse_error":  "",
	}

	for k, v := range t.Templates {
		templates[k] = v
	}

	resp, err := generateOutput("default", addTemplates(t.Template, templates), t.CommonTemplate, t.UseRawOutput)
	if err != nil {
		return "", err
	}

	return resp, nil
}

// SetValue sets template entities to CommonTemplate
func (t *Template) SetValue(ct CommonTemplate) {
	t.CommonTemplate = ct
}

func addTemplates(tpl string, templates map[string]string) string {
	for k, v := range templates {
		tpl += `{{define "` + k + `"}}` + v + "{{end}}"
	}
	return tpl
}

func (t *Template) IsSamePlan(executedStr string) bool {
	templateSplitted := strings.Split(t.Template, "\n")
	planTitleLineIndex := -1
	for i, ts := range templateSplitted {
		if strings.Contains(ts, `template "plan_title"`) {
			planTitleLineIndex = i
		}
	}

	if planTitleLineIndex == -1 {
		return false
	}

	targetSplitted := strings.Split(executedStr, "\n")
	if len(targetSplitted) < planTitleLineIndex+1 {
		return false
	}
	oldTitle := targetSplitted[planTitleLineIndex]

	// attempt to detect changes in the exit code.
	for _, exitCode := range []int{0, 1} {
		commonTemplate := CommonTemplate{
			ExitCode: exitCode,
			Vars:     t.Vars,
		}
		newTitle, err := generateOutput("default", planTitleTemplate, commonTemplate, t.UseRawOutput)
		if err != nil {
			return false
		}

		if newTitle == oldTitle {
			return true
		}
	}

	return false
}
