# tfcmt-gitlab

[![Build Status](https://github.com/hirosassa/tfcmt-gitlab/workflows/test/badge.svg)](https://github.com/hirosassa/tfcmt-gitlab/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/hirosassa/tfcmt-gitlab)](https://goreportcard.com/report/github.com/hirosassa/tfcmt-gitlab)
[![GitHub last commit](https://img.shields.io/github/last-commit/hirosassa/tfcmt-gitlab.svg)](https://github.com/hirosassa/tfcmt-gitlab)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/hirosassa/tfcmt-gitlab/main/LICENSE)

Fork of [suzuki-shunsuke/tfcmt](https://github.com/suzuki-shunsuke/tfcmt), supporting GitLab (dropped GitHub support).

## Document

### Prerequisites

Create and store GitLab access token in [project or group CI variables](https://docs.gitlab.com/ee/ci/variables/#add-a-cicd-variable-to-a-project) with key name `GITLAB_TOKEN`.

ref: [Project access tokens | GitLab](https://docs.gitlab.com/ee/user/project/settings/project_access_tokens.html)


Basic commands are as follows:

```shell
# plan
tfcmt-gitlab plan --patch -- terraform plan -no-color

# apply
tfcmt-gitlab apply -- terraform apply -auto-approve -no-color
```

`tfcmt-gitlab` runs without any configuration file.
The concrete examples of configuration of `tfcmt-gitlab` running on GitLab CI are available in [examples/getting-started](https://github.com/hirosassa/tfcmt-gitlab/tree/main/examples/getting-started).

## License

### License of original code

This is a fork of [mercari/tfnotify](https://github.com/mercari/tfnotify) and [suzuki-shunsuke/tfcmt](https://github.com/suzuki-shunsuke/tfcmt), so about the origincal license, please see https://github.com/mercari/tfnotify#license and https://github.com/suzuki-shunsuke/tfcmt#license.

Copyright 2018 Mercari, Inc.

Licensed under the MIT License.

### License of code which we wrote

MIT
