.PHONY: before-commit
before-commit: check-code
before-commit: test-code

.PHONY: check-code
check-code:
	go fmt ./... && go vet ./...

.PHONY: test-code
test-code:
	go clean -testcache
	go test -v ./...
