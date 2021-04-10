ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")
APP_NAME="golang-metrics"
EXPORTED_PACKAGE_COMMENT_WARNING="exported (var|function|method|type|const) \S+ should have comment"
UNDERSCORED_PACKAGE_NAME_WARNING="don't use an underscore in package name"

build: compile fmt vet lint

compile:
	go build

fmt:
	go fmt $(ALL_PACKAGES)

vet:
	go vet $(ALL_PACKAGES)

lint:
	@for p in $(ALL_PACKAGES); do \
		golint $$p | { grep -vwE $(EXPORTED_PACKAGE_COMMENT_WARNING) || true; } | { grep -vwE $(UNDERSCORED_PACKAGE_NAME_WARNING) || true; } \
	done