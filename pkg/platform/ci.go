package platform

import (
	"fmt"
	"os"
	"strconv"

	"github.com/hirosassa/tfcmt-gitlab/pkg/config"
)

const GitlabCI = "gitlabci"

func Complement(cfg *config.Config) error {
	if err := complementWithCIEnv(&cfg.CI); err != nil {
		return err
	}

	return complementWithGeneric(cfg)
}

func complementWithCIEnv(ci *config.CI) error {
	ci.Name = GitlabCI

	if ci.NameSpace == "" {
		ci.NameSpace = os.Getenv("CI_PROJECT_NAMESPACE")
	}

	if ci.Project == "" {
		ci.Project = os.Getenv("CI_PROJECT_NAME")
	}

	if ci.SHA == "" {
		ci.SHA = os.Getenv("CI_COMMIT_SHA")
	}

	if ci.MRNumber <= 0 {
		mr := os.Getenv("CI_MERGE_REQUEST_IID")
		a, err := strconv.Atoi(mr)
		if err != nil {
			return fmt.Errorf("parse CI_MERGE_REQUEST_IID %s: %w", mr, err)
		}
		ci.MRNumber = a
	}

	if ci.Link == "" {
		ci.Link = os.Getenv("CI_JOB_URL")
	}
	return nil
}

func complementWithGeneric(cfg *config.Config) error {
	gen := generic{
		param: Param{
			NameSpace: cfg.Complement.NameSpace,
			Project:   cfg.Complement.Project,
			SHA:       cfg.Complement.SHA,
			MRNumber:  cfg.Complement.MR,
			Link:      cfg.Complement.Link,
			Vars:      cfg.Complement.Vars,
		},
	}

	if cfg.CI.NameSpace == "" {
		cfg.CI.NameSpace = gen.NameSpace()
	}

	if cfg.CI.Project == "" {
		cfg.CI.Project = gen.Project()
	}

	if cfg.CI.SHA == "" {
		cfg.CI.SHA = gen.SHA()
	}

	if cfg.CI.MRNumber <= 0 {
		n, err := gen.PRNumber()
		if err != nil {
			return err
		}
		cfg.CI.MRNumber = n
	}

	if cfg.CI.Link == "" {
		cfg.CI.Link = gen.Link()
	}

	vars := gen.Vars()
	for k, v := range cfg.Vars {
		vars[k] = v
	}
	cfg.Vars = vars

	return nil
}
