.PHONY: before-commit
before-commit: check test

.PHONY: check
check:
	go fmt ./... && go vet ./...

.PHONY: test
test:
	go clean -testcache
	go test -v ./...

.PHONE: coverage
coverage:
	go clean -cache
	go test -v -coverprofile=c.out ./...
	go tool cover -html=c.out
