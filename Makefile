.PHONY: mockgen
mockgen:
	go run go.uber.org/mock/mockgen -source=./pkg/notifier/gitlab/gitlab.go -destination=./pkg/notifier/gitlab/gen/gitlab.go -package gitlabmock
