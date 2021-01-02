OTELCOL_BUILDER_VERSION ?= 0.2.0

build: build-agent build-collector

build-agent: otelcol-builder
	@mkdir -p builds/agent
	@$(OTELCOL_BUILDER) --config manifests/agent.yaml

build-collector: otelcol-builder
	@mkdir -p builds/collector
	@$(OTELCOL_BUILDER) --config manifests/collector.yaml

otelcol-builder:
ifeq (, $(shell which opentelemetry-collector-builder))
	go get github.com/observatorium/opentelemetry-collector-builder@v$(OTELCOL_BUILDER_VERSION)
endif
OTELCOL_BUILDER=$(shell which opentelemetry-collector-builder)
