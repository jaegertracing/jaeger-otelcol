OTELCOL_BUILDER_VERSION ?= 0.2.0
GOFMT=gofmt
GOLINT=golint
LINT_LOG=.lint.log

# all .go files that are not auto-generated and should be auto-formatted and linted.
ALL_SRC := $(shell find . -type d \( -name builds \) -prune -false -o \
				   -name '*.go' \
				   -not -name '.*' \
				   -type f | \
				sort)

# ALL_PKGS is used with 'golint'
ALL_PKGS := $(shell echo $(dir $(ALL_SRC)) | tr ' ' '\n' | sort -u)

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

.PHONY: e2e-tests
e2e-tests: build e2e-tests-agent-smoke e2e-tests-collector-smoke

.PHONY: e2e-tests-agent-smoke
e2e-tests-agent-smoke: build-agent
	@echo Running Agent end-to-end tests...
	@go test -tags=agent_smoke ./test/e2e/agent/... $(TEST_OPTIONS)
OTELCOL_BUILDER=$(shell which opentelemetry-collector-builder)

.PHONY: e2e-tests-collector-smoke
e2e-tests-collector-smoke: build-collector
	@echo Running Collector end-to-end tests...
	@go test -tags=collector_smoke ./test/e2e/collector/... $(TEST_OPTIONS)

.PHONY: lint
lint: fmt go-lint

.PHONY: go-lint
go-lint:
	@cat /dev/null > $(LINT_LOG)
	@echo Running go lint...
	@$(GOLINT) $(ALL_PKGS) | grep -v _nolint.go >> $(LINT_LOG) || true;
	@[ ! -s "$(LINT_LOG)" ] || (echo "Lint Failures" | cat - $(LINT_LOG) && false)

.PHONY: fmt
fmt:
	@echo Running go fmt on ALL_SRC ...
	@$(GOFMT) -e -s -l -w $(ALL_SRC)

.PHONY: install-tools
install-tools:
	go install golang.org/x/lint/golint
