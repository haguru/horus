.PHONY: build test unittest lint clean update fmt docker run

ARCH=$(shell uname -m)

MICROSERVICE=gateway
DOCKERDATABASE=mongodb/mongodb-community-server:latest

.PHONY: build test clean fmt docker run

APPVERSION=$(shell cat ./VERSION 2>/dev/null || echo 0.0.0)
GIT_SHA=$(shell git rev-parse HEAD)

build:
	CGO_ENABLED=0 GOOS=linux go build -tags "$(ADD_BUILD_TAGS)" $(GOFLAGS) -o $(MICROSERVICE)

tidy:
	go mod tidy

t:
	[ -z "$$(gofmt -p -l . || echo 'err')" ]

unittest:
	go test ./... -coverprofile=coverage.out ./...

lint:
	@which golangci-lint >/dev/null || echo "WARNING: go linter not installed. To install, run make install-lint"
	@if [ "z${ARCH}" = "zx86_64" ] && which golangci-lint >/dev/null ; then golangci-lint run --config .golangci.yml ; else echo "WARNING: Linting skipped (not on x86_64 or linter not installed)"; fi

install-lint:
	sudo curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.54.2

test: unittest lint
	go vet ./...
	gofmt -l $$(find . -type f -name '*.go'| grep -v "/vendor/")
	[ "`gofmt -l $$(find . -type f -name '*.go'| grep -v "/vendor/")`" = "" ]

clean:
	rm -f $(MICROSERVICE)

fmt:
	go fmt ./...

docker:
	docker build \
		--rm \
		-t haguru/follower:$(GIT_SHA) \
		-t haguru/follower:$(APPVERSION)-dev \
		.

run: build
	./$(MICROSERVICE)

launchdb: 
	docker run --rm -d -p 50056:27017/tcp  --name followerdb $(DOCKERDATABASE) 

vendor:
	go mod vendor