package platform

import (
	"fmt"
	"strconv"

	"github.com/hirosassa/tfcmt-gitlab/pkg/domain"
)

type Param struct {
	NameSpace []domain.ComplementEntry
	Project   []domain.ComplementEntry
	SHA       []domain.ComplementEntry
	MRNumber  []domain.ComplementEntry
	Link      []domain.ComplementEntry
	Vars      map[string][]domain.ComplementEntry
}

type generic struct {
	param Param
}

func (gen *generic) render(entries []domain.ComplementEntry) (string, error) {
	var e error
	for _, entry := range entries {
		a, err := entry.Entry()
		if err != nil {
			e = err
			continue
		}
		if a != "" {
			return a, nil
		}
	}
	return "", e
}

func (gen *generic) returnString(entries []domain.ComplementEntry) string {
	s, err := gen.render(entries)
	if err != nil {
		return ""
	}
	return s
}

func (gen *generic) NameSpace() string {
	return gen.returnString(gen.param.NameSpace)
}

func (gen *generic) Project() string {
	return gen.returnString(gen.param.Project)
}

func (gen *generic) SHA() string {
	return gen.returnString(gen.param.SHA)
}

func (gen *generic) Link() string {
	return gen.returnString(gen.param.Link)
}

func (gen *generic) IsMR() bool {
	return gen.returnString(gen.param.MRNumber) != ""
}

func (gen *generic) PRNumber() (int, error) {
	s, err := gen.render(gen.param.MRNumber)
	if err != nil {
		return 0, err
	}
	if s == "" {
		return 0, nil
	}
	b, err := strconv.Atoi(s)
	if err == nil {
		return b, nil
	}
	return 0, fmt.Errorf("parse pull request number as int: %w", err)
}

func (gen *generic) Vars() map[string]string {
	vars := make(map[string]string, len(gen.param.Vars))
	for k, v := range gen.param.Vars {
		a := gen.returnString(v)
		vars[k] = a
	}
	return vars
}
