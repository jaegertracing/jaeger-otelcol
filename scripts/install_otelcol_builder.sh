#!/usr/bin/env bash

while getopts d:v: flag
do
    case "${flag}" in
        d) otelcol_builder_dir=${OPTARG};;
        v) otelcol_builder_version=${OPTARG};;
    esac
done

if [ "$(which opentelemetry-collector-builder)" ]; then
  echo "opentelemetry-collector-builder already installed"
  exit 0
fi

echo "installing opentelemetry-collector-builder"
otelcol_builder="$otelcol_builder_dir/opentelemetry-collector-builder"
goos=$(go env GOOS)
goarch=$(go env GOARCH)

set -ex
mkdir -p "$otelcol_builder_dir"
curl -sLo "$otelcol_builder" "https://github.com/open-telemetry/opentelemetry-collector-builder/releases/download/v${otelcol_builder_version}/opentelemetry-collector-builder_${otelcol_builder_version}_${goos}_${goarch}"
chmod +x "${otelcol_builder}"