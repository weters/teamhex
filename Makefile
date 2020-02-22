all: .git/hooks/pre-commit

.git/hooks/pre-commit:
	ln -s ../../git-hooks/pre-commit .git/hooks/pre-commit

.PHONY: format
format:
	find . \! -path './.git*' -type f -name '*.go' | xargs -L 1 gofmt -s -w

.PHONY: test
test:
	go test -coverprofile=coverage.out ./...

.PHONY: clean
clean:
	rm -f coverage.out
