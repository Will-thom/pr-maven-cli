.PHONY: test test-race coverage coverage-check ci

test:
	go test ./...

test-race:
	go test -race ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

coverage-check:
	PRMAVEN_COVERAGE=1 PRMAVEN_MIN_COVERAGE=70 sh scripts/test.sh

ci: test test-race coverage-check
