package platform

import (
	"fmt"
	"os"
	"strconv"

	"github.com/hirosassa/tfcmt-gitlab/pkg/config"
)

func Complement(cfg *config.Config) error {
	if err := complementCIInfo(&cfg.CI); err != nil {
		return err
	}

	return complementWithGeneric(cfg)
}

func complementCIInfo(ci *config.CI) error {
	if ci.MRNumber <= 0 {
		if mrS := os.Getenv("CI_INFO_MR_NUMBER"); mrS != "" {
			a, err := strconv.Atoi(mrS)
			if err != nil {
				return fmt.Errorf("parse CI_INFO_PR_NUMBER %s: %w", mrS, err)
			}
			ci.MRNumber = a
		}
	}
	return nil
}

func getLink(ciname string) string {
	switch ciname {
	case "gilabci", "gitlab-ci":
		return os.Getenv("CI_JOB_URL")
	}
	return ""
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
