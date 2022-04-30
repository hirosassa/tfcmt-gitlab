package cli

import (
	"github.com/urfave/cli/v2"
)

type LDFlags struct {
	Version string
	Commit  string
	Date    string
}

func (flags *LDFlags) AppVersion() string {
	return flags.Version + " (" + flags.Commit + ")"
}

func New(flags *LDFlags) *cli.App {
	app := cli.NewApp()
	app.Name = "tfcmt-gitlab"
	app.Usage = "Notify the execution result of terraform command"
	app.Version = flags.AppVersion()
	app.Flags = []cli.Flag{
		&cli.StringFlag{Name: "namespace", Usage: "GitLab namespace name"},
		&cli.StringFlag{Name: "project", Usage: "GitLab project name"},
		&cli.StringFlag{Name: "sha", Usage: "commit SHA (revision)"},
		&cli.StringFlag{Name: "build-url", Usage: "build url"},
		&cli.StringFlag{Name: "log-level", Usage: "log level"},
		&cli.IntFlag{Name: "mr", Usage: "merge request number"},
		&cli.StringFlag{Name: "config", Usage: "config path"},
		&cli.StringSliceFlag{Name: "var", Usage: "template variables. The format of value is '<name>:<value>'"},
	}
	app.Commands = []*cli.Command{
		{
			Name:   "plan",
			Usage:  "Run terraform plan and post a comment to GitHub commit or pull request",
			Action: cmdPlan,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "patch",
					Usage: "update an existing comment instead of creating a new comment. If there is no existing comment, a new comment is created.",
				},
			},
		},
		{
			Name:   "apply",
			Usage:  "Run terraform apply and post a comment to GitHub commit or pull request",
			Action: cmdApply,
		},
		{
			Name:  "version",
			Usage: "Show version",
			Action: func(ctx *cli.Context) error {
				cli.ShowVersion(ctx)
				return nil
			},
		},
	}
	return app
}
