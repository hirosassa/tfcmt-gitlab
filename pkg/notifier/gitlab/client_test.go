package gitlab

import (
	"os"
	"testing"
)

func TestNewClient(t *testing.T) { //nolint:paralleltest
	t.Setenv(EnvToken, "")

	testCases := []struct {
		config   Config
		envToken string
		expect   string
	}{
		{
			// specify directly
			config:   Config{Token: "abcdefg"},
			envToken: "",
			expect:   "",
		},
		{
			// specify via env but not to be set env (part 1)
			config:   Config{Token: "GITLAB_TOKEN"},
			envToken: "",
			expect:   "gitlab token is missing",
		},
		{
			// specify via env (part 1)
			config:   Config{Token: "GITLAB_TOKEN"},
			envToken: "abcdefg",
			expect:   "",
		},
		{
			// specify via env but not to be set env (part 2)
			config:   Config{Token: "$GITLAB_TOKEN"},
			envToken: "",
			expect:   "gitlab token is missing",
		},
		{
			// specify via env (part 2)
			config:   Config{Token: "$GITLAB_TOKEN"},
			envToken: "abcdefg",
			expect:   "",
		},
		{
			// no specification (part 1)
			config:   Config{},
			envToken: "",
			expect:   "gitlab token is missing",
		},
		{
			// no specification (part 2)
			config:   Config{},
			envToken: "abcdefg",
			expect:   "gitlab token is missing",
		},
	}
	for _, testCase := range testCases {
		os.Setenv(EnvToken, testCase.envToken)
		_, err := NewClient(testCase.config)
		if err == nil {
			continue
		}
		if err.Error() != testCase.expect {
			t.Errorf("got %q but want %q", err.Error(), testCase.expect)
		}
	}
}

func TestNewClientWithBaseURL(t *testing.T) { //nolint:paralleltest
	t.Setenv(EnvBaseURL, "")

	testCases := []struct {
		config     Config
		envBaseURL string
		expect     string
	}{
		{
			// specify directly
			config: Config{
				Token:   "abcdefg",
				BaseURL: "https://git.example.com/",
			},
			envBaseURL: "",
			expect:     "https://git.example.com/api/v4/",
		},
		{
			// specify via env but not to be set env (part 1)
			config: Config{
				Token:   "abcdefg",
				BaseURL: "GITLAB_BASE_URL",
			},
			envBaseURL: "",
			expect:     "https://gitlab.com/api/v4/",
		},
		{
			// specify via env (part 1)
			config: Config{
				Token:   "abcdefg",
				BaseURL: "GITLAB_BASE_URL",
			},
			envBaseURL: "https://git.example.com/",
			expect:     "https://git.example.com/api/v4/",
		},
		{
			// specify via env but not to be set env (part 2)
			config: Config{
				Token:   "abcdefg",
				BaseURL: "$GITLAB_BASE_URL",
			},
			envBaseURL: "",
			expect:     "https://gitlab.com/api/v4/",
		},
		{
			// specify via env (part 2)
			config: Config{
				Token:   "abcdefg",
				BaseURL: "$GITLAB_BASE_URL",
			},
			envBaseURL: "https://git.example.com/",
			expect:     "https://git.example.com/api/v4/",
		},
		{
			// no specification (part 1)
			config:     Config{Token: "abcdefg"},
			envBaseURL: "",
			expect:     "https://gitlab.com/api/v4/",
		},
		{
			// no specification (part 2)
			config:     Config{Token: "abcdefg"},
			envBaseURL: "https://git.example.com/",
			expect:     "https://gitlab.com/api/v4/",
		},
	}
	for _, testCase := range testCases {
		os.Setenv(EnvBaseURL, testCase.envBaseURL)
		c, err := NewClient(testCase.config)
		if err != nil {
			continue
		}
		url := c.Client.BaseURL().String()
		if url != testCase.expect {
			t.Errorf("got %q but want %q", url, testCase.expect)
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
