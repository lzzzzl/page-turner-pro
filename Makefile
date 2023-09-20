all: build

GIT_BRANCH=$(shell git branch | grep \* | cut -d ' ' -f2)
GIT_REV=$(shell git rev-parse HEAD | cut -c1-7)
BUILD_DATE=$(shell date +%Y-%m-%d.%H:%M:%S)
EXTRA_LD_FLAGS=-X main.AppVersion=${GIT_BRANCH}-${GIT_REV} -X main.AppBuild=${BUILD_DATE}

GOLANGCI_LINT_VERSION="v1.42.1"
DATABASE_DSN="postgresql://page_turner_pro:page_turner_pro@localhost:5432/page_turner_pro?sslmode=disable"
TEST_PACKAGES=./internal/...

build:
	go build -ldflags '${EXTRA_LD_FLAGS}' -o bin/page-turner-pro ./cmd/appserver

run: build
		./bin/page-turner-pro \
		--database_dsn=${DATABASE_DSN} \
		| jq

test:
	go vet $(TEST_PACKAGES)
	go test -cover -coverprofile cover.out $(TEST_PACKAGES)
	go tool cover -func=cover.out | tail -n 1
