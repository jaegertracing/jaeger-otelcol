module github.com/jaegertracing/jaeger-otelcol

go 1.15

require (
	github.com/jaegertracing/jaeger v1.21.1-0.20210208225804-d23b3e234ed2
	github.com/jaegertracing/jaeger-otelcol/test/tools/tracegen v0.0.0-20201218111612-2f0546350989
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/prom2json v1.3.0
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/collector v0.20.0
	go.uber.org/zap v1.16.0
	golang.org/x/lint v0.0.0-20201208152925-83fdc39ff7b5 // indirect
)

replace honnef.co/go/tools => honnef.co/go/tools v0.0.1-2020.1.6
