// Copyright (c) 2020 The Jaeger Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package e2e

import (
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"

	"go.uber.org/zap"

	pcm "github.com/prometheus/client_model/go"
	"github.com/prometheus/prom2json"
	"github.com/stretchr/testify/require"
)

var (
	// LogLevel is used to set the level for the zap logger.
	LogLevel = getStringEnv("LOG_LEVEL", "info")
)

func getStringEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

// StartCollector starts the executable in background.
func StartCollector(t *testing.T, logger zap.SugaredLogger, executable, configFileName string, loggerOutput io.Writer, metricsPort string) *exec.Cmd {
	// metrics-addr is needed to avoid collisions when we start multiple processes
	arguments := []string{executable, "--config", configFileName, "--metrics-addr", "localhost:" + metricsPort}

	cmd := exec.Command(executable, arguments...)
	cmd.Stderr = loggerOutput
	err := cmd.Start()
	require.NoError(t, err)

	logger.Infof("Started process %d with %s", cmd.Process.Pid, executable)
	return cmd
}

// GetPrometheusCounter returns the counter named from the specified endpoint.
func GetPrometheusCounter(t *testing.T, metricsEndpoint, metricName string) float64 {
	counter := GetPrometheusMetric(t, metricsEndpoint, metricName)
	return *counter.Metric[0].Counter.Value
}

// GetPrometheusMetric returns the metric named from the specified endpoint.
func GetPrometheusMetric(t *testing.T, metricsEndpoint, metricName string) pcm.MetricFamily {
	allMetrics := GetPrometheusMetrics(t, metricsEndpoint)
	return allMetrics[metricName]
}

// GetPrometheusMetrics returns all metrics from the specified endpoint.
func GetPrometheusMetrics(t *testing.T, metricsEndpoint string) map[string]pcm.MetricFamily {
	// This code is mostly copied from https://github.com/prometheus/prom2json except it
	// returns MetricFamily objects as that is more useful than JSON for tests.
	mfChan := make(chan *pcm.MetricFamily, 1024)
	err := prom2json.FetchMetricFamilies(metricsEndpoint, mfChan, &http.Transport{})
	require.NoError(t, err)
	result := map[string]pcm.MetricFamily{}
	for mf := range mfChan {
		result[*mf.Name] = *mf
	}

	return result
}

// CreateTempFile creates a temp file.
func CreateTempFile(t *testing.T) *os.File {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "jaeger-otel-test-")
	require.NoError(t, err)
	return tmpFile
}

// GetFreePort will return a free tcp port.
func GetFreePort(t *testing.T) string {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	require.NoError(t, err)
	listener, err := net.ListenTCP("tcp", addr)
	require.NoError(t, err)
	defer listener.Close()

	address := listener.Addr().String()
	colon := strings.Index(address, ":")
	port := address[colon+1:]
	return port
}

// GetLogger returns a Zap Sugared logger.
func GetLogger(t *testing.T) zap.SugaredLogger {
	var zapLogger *zap.Logger
	var err error
	if LogLevel == "info" {
		zapLogger, err = zap.NewProduction()
	} else {
		zapLogger, err = zap.NewDevelopment()
	}
	require.NoError(t, err)
	return *zapLogger.Sugar()
}
