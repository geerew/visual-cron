# A helper makefile 
#
# Requires 
#  * make (yum, apt, chocolatey)
#  * go (http://golang.org) 
#  * gox (go get github.com/mitchellh/gox)

DISTROS="linux/amd64 darwin/amd64 windows/amd64"
PACKAGES := $(shell go list ./... | grep -v /vendor/)

.PHONY: test
test: # run unit tests
	@go clean -testcache
	@$(foreach pkg,$(PACKAGES), \
		go test -p=1 -cover -covermode=count -coverprofile=coverage.out ${pkg})

.PHONY: test-cover
test-cover: test # run unit tests and show test coverage information
	go tool cover -html="coverage.out"

.PHONY: build
build: test # build the binary
	@gox -verbose -osarch ${DISTROS} -output "builds/visualcron_{{.OS}}_{{.Arch}}"

.PHONY: run
run: # run the application
	@go run .

.PHONY: fmt
fmt: # run "go fmt" on all Go packages
	@go fmt $(PACKAGES)