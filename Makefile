.ONESHELL:
SHELL := /bin/bash

TEST_COUNT?=1
ACCTEST_PARALLELISM?=4
ACCTEST_TIMEOUT?=10m

all: test testacc build

tidy:
	go mod tidy

fmt:
	go fmt ./...

vet:
	go vet ./...

test: tidy fmt vet
	go test ./... -coverprofile cover.out

testacc: tidy fmt vet
	TF_ACC=1 go test ./pkg/provider -v -count $(TEST_COUNT) -parallel $(ACCTEST_PARALLELISM) -timeout $(ACCTEST_TIMEOUT)

build:
	CGO_ENABLED=0 go build -o ./bin/pal main.go

install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/xenitab/pal/0.0.0-dev/$${GOOS}_$${GOARCH}
	cp ./bin/pal ~/.terraform.d/plugins/registry.terraform.io/xenitab/pal/0.0.0-dev/$${GOOS}_$${GOARCH}/terraform-provider-pal

.PHONY: docs
docs: tools
	tfplugindocs generate

tools:
	GO111MODULE=on go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

.SILENT:
lint:
	set -e
	echo lint: Start
	EXAMPLES=$$(find examples -mindepth 1 -maxdepth 1 -type d)
	DELETE=examples/data-sources
	echo $${array[@]/$$DELETE}
	EXAMPLES=( "$${EXAMPLES[@]/$$DELETE}" )
	for EXAMPLE in $${EXAMPLES}; do
					echo $$EXAMPLE
					tflint -c examples/.tflint.hcl $${EXAMPLE}
	done
	echo lint: Success
