package gitlab

import (
	"github.com/hirosassa/tfcmt-gitlab/pkg/notifier"
	"github.com/hirosassa/tfcmt-gitlab/pkg/terraform"
	"github.com/sirupsen/logrus"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// NotifyService handles communication with the notification related
// methods of GitLab API
type NotifyService service

// Notify posts comment optimized for notifications
func (g *NotifyService) Notify(param notifier.ParamExec) (int, error) { //nolint:cyclop
	cfg := g.client.Config
	parser := g.client.Config.Parser
	template := g.client.Config.Template
	var errMsgs []string

	result := parser.Parse(param.CombinedOutput)
	result.ExitCode = param.ExitCode
	if result.HasParseError {
		template = g.client.Config.ParseErrorTemplate
	} else {
		if result.Error != nil {
			return result.ExitCode, result.Error
		}
		if result.Result == "" {
			return result.ExitCode, result.Error
		}
	}

	_, isPlan := parser.(*terraform.PlanParser)
	if isPlan {
		if cfg.MR.IsNumber() && cfg.ResultLabels.HasAnyLabelDefined() {
			errMsgs = append(errMsgs, g.updateLabels(result)...)
		}
	}

	template.SetValue(terraform.CommonTemplate{
		Result:                 result.Result,
		ChangedResult:          result.ChangedResult,
		ChangeOutsideTerraform: result.OutsideTerraform,
		Warning:                result.Warning,
		Link:                   cfg.CI,
		UseRawOutput:           cfg.UseRawOutput,
		HasDestroy:             result.HasDestroy,
		Vars:                   cfg.Vars,
		Templates:              cfg.Templates,
		Stdout:                 param.Stdout,
		Stderr:                 param.Stderr,
		CombinedOutput:         param.CombinedOutput,
		ExitCode:               param.ExitCode,
		ErrorMessages:          errMsgs,
		CreatedResources:       result.CreatedResources,
		UpdatedResources:       result.UpdatedResources,
		DeletedResources:       result.DeletedResources,
		ReplacedResources:      result.ReplacedResources,
	})
	body, err := template.Execute()
	if err != nil {
		return result.ExitCode, err
	}

	_, isApply := parser.(*terraform.ApplyParser)
	logE := logrus.WithFields(logrus.Fields{
		"program": "tfcmt",
	})

	if !isApply && cfg.Patch && cfg.MR.Number != 0 {
		logE.Debug("try patching")
		// If fail to list comments, try to create new post.
		comments, _ := g.client.Comment.List(cfg.MR.Number)
		for i := len(comments) - 1; i >= 0; i-- {
			comment := comments[i]

			if template.IsSamePlan(comment.Body) {
				logE.Debugf("Patch comment from `%s` to `%s`", comment.Body, body)
				if err := g.client.Comment.Patch(int(comment.ID), body, PostOptions{
					Number:   cfg.MR.Number,
					Revision: cfg.MR.Revision,
				}); err != nil {
					return result.ExitCode, err
				}
				return result.ExitCode, nil
			}
		}
		logE.WithField("size", len(comments)).Debug("list comments")
	}

	if result.HasNoChanges && result.Warning == "" && len(errMsgs) == 0 && cfg.SkipNoChanges {
		logE.Debug("skip posting a comment because there is no change")
		return result.ExitCode, nil
	}

	logE.Debug("create a comment")

	if err := g.client.Comment.Post(body, PostOptions{
		Number:   cfg.MR.Number,
		Revision: cfg.MR.Revision,
	}); err != nil {
		return result.ExitCode, err
	}
	return result.ExitCode, nil
}

func (g *NotifyService) updateLabels(result terraform.ParseResult) []string { //nolint:cyclop
	cfg := g.client.Config
	var (
		labelToAdd string
		labelColor string
	)

	switch {
	case result.HasAddOrUpdateOnly:
		labelToAdd = cfg.ResultLabels.AddOrUpdateLabel
		labelColor = cfg.ResultLabels.AddOrUpdateLabelColor
	case result.HasDestroy:
		labelToAdd = cfg.ResultLabels.DestroyLabel
		labelColor = cfg.ResultLabels.DestroyLabelColor
	case result.HasNoChanges:
		labelToAdd = cfg.ResultLabels.NoChangesLabel
		labelColor = cfg.ResultLabels.NoChangesLabelColor
	case result.HasPlanError:
		labelToAdd = cfg.ResultLabels.PlanErrorLabel
		labelColor = cfg.ResultLabels.PlanErrorLabelColor
	}

	errMsgs := []string{}

	logE := logrus.WithFields(logrus.Fields{
		"program": "tfcmt",
	})

	currentLabelColor, err := g.removeResultLabels(labelToAdd)
	if err != nil {
		msg := "remove labels: " + err.Error()
		logE.WithError(err).Error("remove labels")
		errMsgs = append(errMsgs, msg)
	}

	if labelToAdd == "" {
		return errMsgs
	}

	if currentLabelColor == "" {
		labels, err := g.client.API.AddMergeRequestLabels(&[]string{labelToAdd}, cfg.MR.Number)
		if err != nil {
			msg := "add a label " + labelToAdd + ": " + err.Error()
			logE.WithError(err).WithFields(logrus.Fields{
				"label": labelToAdd,
			}).Error("add a label")
			errMsgs = append(errMsgs, msg)
		}
		if labelColor != "" {
			// set the color of label
			for _, label := range labels {
				if labelToAdd == label {
					l, _, err := g.client.API.GetLabel(label)
					if err != nil {
						msg := "failed to get Label " + label + ": " + err.Error()
						logE.WithError(err).WithFields(logrus.Fields{
							"label": labelToAdd,
						}).Error("get a label")
						errMsgs = append(errMsgs, msg)
					}

					if l.Color != labelColor {
						if _, _, err := g.client.API.UpdateLabel(&gitlab.UpdateLabelOptions{Name: &labelToAdd, Color: &labelColor}); err != nil {
							msg := "update a label color (name: " + labelToAdd + ", color: " + labelColor + "): " + err.Error()
							logE.WithError(err).WithFields(logrus.Fields{
								"label": labelToAdd,
								"color": labelColor,
							}).Error("update a label color")
							errMsgs = append(errMsgs, msg)
						}
					}
				}
			}
		}
	} else if labelColor != "" && labelColor != currentLabelColor {
		// set the color of label
		if _, _, err := g.client.API.UpdateLabel(&gitlab.UpdateLabelOptions{Name: &labelToAdd, Color: &labelColor}); err != nil {
			msg := "update a label color (name: " + labelToAdd + ", color: " + labelColor + "): " + err.Error()
			logE.WithError(err).WithFields(logrus.Fields{
				"label": labelToAdd,
				"color": labelColor,
			}).Error("update a label color")
			errMsgs = append(errMsgs, msg)
		}
	}
	return errMsgs
}

func (g *NotifyService) removeResultLabels(label string) (string, error) {
	cfg := g.client.Config
	labels, err := g.client.API.ListMergeRequestLabels(cfg.MR.Number, nil)
	if err != nil {
		return "", err
	}

	labelColor := ""
	for _, l := range labels {
		labelText := l
		if labelText == label {
			currentLabel, _, err := g.client.API.GetLabel(l)
			if err != nil {
				return "", err
			}
			labelColor = currentLabel.Color
			continue
		}
		if cfg.ResultLabels.IsResultLabel(labelText) {
			_, err := g.client.API.RemoveMergeRequestLabels(&[]string{labelText}, cfg.MR.Number)
			if err != nil {
				return labelColor, err
			}
		}
	}

	return labelColor, nil
}
