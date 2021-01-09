// +build tools

package jaeger_otelcol

// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
import (
	_ "github.com/observatorium/opentelemetry-collector-builder"
	_ "golang.org/x/lint/golint"
)
