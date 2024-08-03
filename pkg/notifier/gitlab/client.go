package gitlab

import (
	"errors"
	"os"
	"strings"

	"github.com/hirosassa/tfcmt-gitlab/pkg/terraform"
	gitlab "github.com/xanzy/go-gitlab"
)

// EnvToken is GitLab API Token
const EnvToken = "GITLAB_TOKEN"

// EnvBaseURL is GitLab base URL. This can be set to a domain endpoint to use with Private GitLab.
const EnvBaseURL = "GITLAB_BASE_URL"

// Client ...
type Client struct {
	*gitlab.Client
	Debug  bool
	Config Config

	common service

	Comment *CommentService
	Commits *CommitsService
	Notify  *NotifyService

	API API
}

// Config is a configuration for GitLab client
type Config struct {
	Token     string
	BaseURL   string
	NameSpace string
	Project   string
	MR        MergeRequest
	CI        string
	Parser    terraform.Parser
	// Template is used for all Terraform command output
	Template           *terraform.Template
	ParseErrorTemplate *terraform.Template
	// ResultLabels is a set of labels to apply depending on the plan result
	ResultLabels     ResultLabels
	Vars             map[string]string
	EmbeddedVarNames []string
	Templates        map[string]string
	UseRawOutput     bool
	Patch            bool
	SkipNoChanges    bool
}

// MergeRequest represents GitLab Merge Request metadata
type MergeRequest struct {
	Revision string
	Title    string
	Message  string
	Number   int
}

type service struct {
	client *Client
}

// NewClient returns Client initialized with Config
func NewClient(cfg Config) (*Client, error) {
	token := getToken(cfg)
	if token == "" {
		return &Client{}, errors.New("gitlab token is missing")
	}

	client, err := gitlab.NewClient(token)
	if err != nil {
		return &Client{}, errors.New("failed to create a new gitlab api client")
	}

	if baseURL := getBaseURL(cfg); baseURL != "" {
		client, err = gitlab.NewClient(token, gitlab.WithBaseURL(baseURL))
		if err != nil {
			return &Client{}, errors.New("failed to create a new gitlab api client")
		}
	}

	c := &Client{
		Config: cfg,
		Client: client,
	}
	c.common.client = c
	c.Comment = (*CommentService)(&c.common)
	c.Commits = (*CommitsService)(&c.common)
	c.Notify = (*NotifyService)(&c.common)

	c.API = &GitLab{
		Client:    client,
		namespace: cfg.NameSpace,
		project:   cfg.Project,
	}

	return c, nil
}

// IsNumber returns true if MergeRequest is Merge Request build
func (mr *MergeRequest) IsNumber() bool {
	return mr.Number != 0
}

// ResultLabels represents the labels to add to the PR depending on the plan result
type ResultLabels struct {
	AddOrUpdateLabel      string
	DestroyLabel          string
	NoChangesLabel        string
	PlanErrorLabel        string
	AddOrUpdateLabelColor string
	DestroyLabelColor     string
	NoChangesLabelColor   string
	PlanErrorLabelColor   string
}

// HasAnyLabelDefined returns true if any of the internal labels are set
func (r *ResultLabels) HasAnyLabelDefined() bool {
	return r.AddOrUpdateLabel != "" || r.DestroyLabel != "" || r.NoChangesLabel != "" || r.PlanErrorLabel != ""
}

// IsResultLabel returns true if a label matches any of the internal labels
func (r *ResultLabels) IsResultLabel(label string) bool {
	switch label {
	case "":
		return false
	case r.AddOrUpdateLabel, r.DestroyLabel, r.NoChangesLabel, r.PlanErrorLabel:
		return true
	default:
		return false
	}
}

func getToken(cfg Config) string {
	if cfg.Token == EnvToken { // specify via config default
		return os.Getenv(EnvToken)
	}

	if strings.HasPrefix(cfg.Token, "$") { // specify via env
		return os.Getenv(strings.TrimPrefix(cfg.Token, "$"))
	}

	if cfg.Token != "" && !strings.HasPrefix(cfg.Token, "$") { // sepcify directly
		return cfg.Token
	}

	return os.Getenv(EnvToken) // specify via default
}

func getBaseURL(cfg Config) string {
	if cfg.BaseURL == EnvBaseURL { // specify via config default
		return os.Getenv(EnvBaseURL)
	}

	if strings.HasPrefix(cfg.BaseURL, "$") { // specify via env
		return os.Getenv(strings.TrimPrefix(cfg.BaseURL, "$"))
	}

	if cfg.BaseURL != "" && !strings.HasPrefix(cfg.BaseURL, "$") { // specify directly
		return cfg.BaseURL
	}

	if os.Getenv("CI_SERVER_URL") != "" { // specify via CI env
		return os.Getenv("CI_SERVER_URL")
	}

	return "" // specify default
}
