default: build

SETENV=
ifeq ($(OS),Windows_NT)
	SETENV=set
endif

lefthook:
	@go install github.com/evilmartians/lefthook@latest
	lefthook install

build:
	go build -v ./...

install: build
	go install -v ./...

lint:
	golangci-lint run

generate:
	go generate ./...

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -tags=all -timeout=120s -parallel=4 ./...

testacc:
	$(SETENV) TF_ACC=1 && go test -v -cover -tags=all -timeout 120m ./...

.PHONY: build install lint generate fmt test testacc
