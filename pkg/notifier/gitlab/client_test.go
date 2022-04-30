package gitlab

import (
	"strings"
	"testing"
)

func TestNewClient(t *testing.T) { //nolint:paralleltest
	testCases := []struct {
		name     string
		config   Config
		envToken string
		expect   string
	}{
		{
			name:     "specify directly",
			config:   Config{Token: "abcdefg"},
			envToken: "",
			expect:   "",
		},
		{
			name:     "specify via env but not to be set env (part 1)",
			config:   Config{Token: "GITLAB_TOKEN"},
			envToken: "",
			expect:   "gitlab token is missing",
		},
		{
			name:     "specify via env (part 1)",
			config:   Config{Token: "GITLAB_TOKEN"},
			envToken: "abcdefg",
			expect:   "",
		},
		{
			name:     "specify via env but not to be set env (part 2)",
			config:   Config{Token: "$GITLAB_TOKEN"},
			envToken: "",
			expect:   "gitlab token is missing",
		},
		{
			name:     "specify via env but not to be set env (part 3)",
			config:   Config{Token: "$TFCMT_GITLAB_TOKEN"},
			envToken: "",
			expect:   "gitlab token is missing",
		},
		{
			name:     "specify via env (part 2)",
			config:   Config{Token: "$GITLAB_TOKEN"},
			envToken: "abcdefg",
			expect:   "",
		},
		{
			name:     "specify via env (part 3)",
			config:   Config{Token: "$TFCMT_GITLAB_TOKEN"},
			envToken: "abcdefg",
			expect:   "",
		},
		{
			name:     "no specification (part 1)",
			config:   Config{},
			envToken: "",
			expect:   "gitlab token is missing",
		},
		{
			name:     "no specification (part 2)",
			config:   Config{},
			envToken: "abcdefg",
			expect:   "gitlab token is missing",
		},
	}
	for _, testCase := range testCases {
		if strings.HasPrefix(testCase.config.Token, "$") || testCase.config.Token == "GITLAB_TOKEN" {
			key := strings.TrimPrefix(testCase.config.Token, "$")
			t.Setenv(key, testCase.envToken)
		}

		_, err := NewClient(testCase.config)
		if err == nil {
			continue
		}
		if err.Error() != testCase.expect {
			t.Errorf("test case %s, got %q but want %q", testCase.name, err.Error(), testCase.expect)
		}
	}
}

func TestNewClientWithBaseURL(t *testing.T) { //nolint:paralleltest
	testCases := []struct {
		name       string
		config     Config
		envBaseURL string
		ciEnvURL   string
		expect     string
	}{
		{
			name: "specify directly",
			config: Config{
				Token:   "abcdefg",
				BaseURL: "https://git.example.com/",
			},
			envBaseURL: "",
			ciEnvURL:   "",
			expect:     "https://git.example.com/api/v4/",
		},
		{
			name: "specify via env but not to be set env (part 1)",
			config: Config{
				Token:   "abcdefg",
				BaseURL: "GITLAB_BASE_URL",
			},
			envBaseURL: "",
			ciEnvURL:   "",
			expect:     "https://gitlab.com/api/v4/",
		},
		{
			name: "specify via env (part 1)",
			config: Config{
				Token:   "abcdefg",
				BaseURL: "GITLAB_BASE_URL",
			},
			envBaseURL: "https://git.example.com/",
			ciEnvURL:   "",
			expect:     "https://git.example.com/api/v4/",
		},
		{
			name: "specify via env but not to be set env (part 2)",
			config: Config{
				Token:   "abcdefg",
				BaseURL: "$GITLAB_BASE_URL",
			},
			envBaseURL: "",
			ciEnvURL:   "",
			expect:     "https://gitlab.com/api/v4/",
		},
		{
			name: "specify via env (part 2)",
			config: Config{
				Token:   "abcdefg",
				BaseURL: "$GITLAB_BASE_URL",
			},
			envBaseURL: "https://git.example.com/",
			ciEnvURL:   "",
			expect:     "https://git.example.com/api/v4/",
		},
		{
			name:       "no specification (part 1)",
			config:     Config{Token: "abcdefg"},
			envBaseURL: "",
			ciEnvURL:   "",
			expect:     "https://gitlab.com/api/v4/",
		},
		{
			name:       "no specification (part 2)",
			config:     Config{Token: "abcdefg"},
			envBaseURL: "https://git.example.com/",
			ciEnvURL:   "",
			expect:     "https://gitlab.com/api/v4/",
		},
		{
			name:       "no specification (part 3)",
			config:     Config{Token: "abcdefg"},
			envBaseURL: "https://git.example.com/",
			ciEnvURL:   "https://gitlab.ci.example.com/",
			expect:     "https://gitlab.ci.example.com/api/v4/",
		},
	}
	for _, testCase := range testCases {
		if strings.HasPrefix(testCase.config.BaseURL, "$") || testCase.config.BaseURL == "GITLAB_BASE_URL" {
			key := strings.TrimPrefix(testCase.config.BaseURL, "$")
			t.Setenv(key, testCase.envBaseURL)
		}
		if testCase.ciEnvURL != "" {
			t.Setenv("CI_SERVER_URL", testCase.ciEnvURL)
		}

		c, err := NewClient(testCase.config)
		if err != nil {
			continue
		}
		url := c.Client.BaseURL().String()
		if url != testCase.expect {
			t.Errorf("test case %s, got %q but want %q", testCase.name, url, testCase.expect)
		}
	}
}

func TestIsNumber(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		mr   MergeRequest
		isMR bool
	}{
		{
			mr: MergeRequest{
				Number: 0,
			},
			isMR: false,
		},
		{
			mr: MergeRequest{
				Number: 123,
			},
			isMR: true,
		},
	}
	for _, testCase := range testCases {
		if testCase.mr.IsNumber() != testCase.isMR {
			t.Errorf("got %v but want %v", testCase.mr.IsNumber(), testCase.isMR)
		}
	}
}
