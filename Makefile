ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")

build: compile fmt vet lint

compile:
	go build

fmt:
	go fmt $(ALL_PACKAGES)

vet:
	go vet $(ALL_PACKAGES)

lint:
	@for p in $(ALL_PACKAGES); do \
		golint $$p | { grep -vwE "exported (var|function|method|type|const) \S+ should have comment" || true; } \
	done