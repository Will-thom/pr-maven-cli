.PHONY: quality test test-race coverage coverage-check build ci

quality:
	sh scripts/quality.sh

test:
	go test ./...

test-race:
	go test -race ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

coverage-check:
	PRMAVEN_COVERAGE=1 PRMAVEN_MIN_COVERAGE=70 sh scripts/test.sh

build:
	sh scripts/build.sh

ci: quality test-race coverage-check build
