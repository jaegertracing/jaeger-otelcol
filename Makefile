OTELCOL_BUILDER_VERSION ?= 0.2.0
GOFMT=gofmt

# all .go files that are not auto-generated and should be auto-formatted and linted.
ALL_SRC := $(shell find . -name '*.go' \
				   -not -name 'doc.go' \
				   -not -name '_*' \
				   -not -name '.*' \
				   -not -name 'gen_assets.go' \
				   -not -name 'mocks*' \
				   -not -name 'model.pb.go' \
				   -not -name 'model_test.pb.go' \
				   -not -name 'storage_test.pb.go' \
				   -not -path './examples/*' \
				   -not -path './vendor/*' \
				   -not -path '*/mocks/*' \
				   -not -path '*/*-gen/*' \
				   -not -path '*/thrift-0.9.2/*' \
				   -type f | \
				sort)

.PHONY: build
build: build-agent build-collector

.PHONY: build-agent
build-agent: otelcol-builder
	@mkdir -p builds/agent
	@$(OTELCOL_BUILDER) --config manifests/agent.yaml

.PHONY: build-collector
build-collector: otelcol-builder
	@mkdir -p builds/collector
	@$(OTELCOL_BUILDER) --config manifests/collector.yaml

.PHONY: otelcol-builder
otelcol-builder:
ifeq (, $(shell which opentelemetry-collector-builder))
	go get github.com/observatorium/opentelemetry-collector-builder@v$(OTELCOL_BUILDER_VERSION)
endif
OTELCOL_BUILDER=$(shell which opentelemetry-collector-builder)

.PHONY: fmt
fmt:
	@echo Running go fmt on ALL_SRC ...
	@$(GOFMT) -e -s -l -w $(ALL_SRC)
