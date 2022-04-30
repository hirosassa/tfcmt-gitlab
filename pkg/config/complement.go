package config

import (
	"errors"

	"github.com/hirosassa/tfcmt-gitlab/pkg/domain"
)

type Complement struct {
	MR        []domain.ComplementEntry
	NameSpace []domain.ComplementEntry
	Project   []domain.ComplementEntry
	SHA       []domain.ComplementEntry
	Link      []domain.ComplementEntry
	Vars      map[string][]domain.ComplementEntry
}

type rawComplement struct {
	MR        []map[string]interface{}
	NameSpace []map[string]interface{}
	Project   []map[string]interface{}
	SHA       []map[string]interface{}
	Link      []map[string]interface{}
	Vars      map[string][]map[string]interface{}
}

func convComplementEntries(maps []map[string]interface{}) ([]domain.ComplementEntry, error) {
	entries := make([]domain.ComplementEntry, len(maps))
	for i, m := range maps {
		entry, err := convComplementEntry(m)
		if err != nil {
			return nil, err
		}
		entries[i] = entry
	}
	return entries, nil
}

func convComplementEntry(m map[string]interface{}) (domain.ComplementEntry, error) {
	t, ok := m["type"]
	if !ok {
		return nil, errors.New(`"type" is required`)
	}
	typ, ok := t.(string)
	if !ok {
		return nil, errors.New(`"type" must be string`)
	}
	switch typ {
	case "envsubst":
		entry := ComplementEnvsubstEntry{}
		if err := newComplementEnvsubstEntry(m, &entry); err != nil {
			return nil, err
		}
		return &entry, nil
	case "template":
		entry := ComplementTemplateEntry{}
		if err := newComplementTemplateEntry(m, &entry); err != nil {
			return nil, err
		}
		return &entry, nil
	default:
		return nil, errors.New(`unsupported type: ` + typ)
	}
}

func (cpl *Complement) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var val rawComplement
	if err := unmarshal(&val); err != nil {
		return err
	}

	mr, err := convComplementEntries(val.MR)
	if err != nil {
		return err
	}
	cpl.MR = mr

	namespace, err := convComplementEntries(val.NameSpace)
	if err != nil {
		return err
	}
	cpl.NameSpace = namespace

	project, err := convComplementEntries(val.Project)
	if err != nil {
		return err
	}
	cpl.Project = project

	sha, err := convComplementEntries(val.SHA)
	if err != nil {
		return err
	}
	cpl.SHA = sha

	link, err := convComplementEntries(val.Link)
	if err != nil {
		return err
	}
	cpl.Link = link

	vars := make(map[string][]domain.ComplementEntry, len(val.Vars))
	for k, v := range val.Vars {
		a, err := convComplementEntries(v)
		if err != nil {
			return err
		}
		vars[k] = a
	}
	cpl.Vars = vars

	return nil
}
