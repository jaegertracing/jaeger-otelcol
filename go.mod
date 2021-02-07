module github.com/jaegertracing/jaeger-otelcol

go 1.15

require (
	github.com/elastic/go-elasticsearch/v6 v6.8.10
	github.com/elastic/go-elasticsearch/v7 v7.0.0
	github.com/jaegertracing/jaeger v1.21.1-0.20210205101436-169322f98c13
	github.com/jaegertracing/jaeger-otelcol/test/tools/tracegen v0.0.0-20201218111612-2f0546350989
	github.com/jaegertracing/jaeger/cmd/opentelemetry v0.0.0-20210205101436-169322f98c13
	github.com/observatorium/opentelemetry-collector-builder v0.4.0
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/prom2json v1.3.0
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/uber/jaeger-lib v2.4.0+incompatible
	go.opencensus.io v0.22.5
	go.opentelemetry.io/collector v0.19.0
	go.uber.org/zap v1.16.0
	golang.org/x/lint v0.0.0-20201208152925-83fdc39ff7b5
	google.golang.org/protobuf v1.25.0 // indirect
)
