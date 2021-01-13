package e2e

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"testing"

	pcm "github.com/prometheus/client_model/go"
	"github.com/prometheus/prom2json"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

var (
	logrusLevel = getStringEnv("LOGRUS_LEVEL", "info")
)

func getStringEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

// StartCollector starts the executable in background
func StartCollector(t *testing.T, executable, configFileName string, loggerOutput io.Writer, metricsPort string) *exec.Cmd {
	// metrics-addr is needed to avoid collisions when we start multiple processes
	arguments := []string{executable, "--config", configFileName, "--metrics-addr", "localhost:" + metricsPort}

	cmd := exec.Command(executable, arguments...)
	cmd.Stderr = loggerOutput
	err := cmd.Start()
	require.NoError(t, err)

	logrus.Infof("Started process %d with %s", cmd.Process.Pid, executable)
	return cmd
}

// GetPrometheusCounter returns the counter named from the specified endpoint
func GetPrometheusCounter(t *testing.T, metricsEndpoint, metricName string) float64 {
	counter := GetPrometheusMetric(t, metricsEndpoint, metricName)
	return *counter.Metric[0].Counter.Value
}

// GetPrometheusMetric returns the metric named from the specified endpoint
func GetPrometheusMetric(t *testing.T, metricsEndpoint, metricName string) pcm.MetricFamily {
	allMetrics := GetPrometheusMetrics(t, metricsEndpoint)
	return allMetrics[metricName]
}

// GetPrometheusMetrics returns all metrics from the specified endpoint
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

// SetLogrusLevel sets the logging level
func SetLogrusLevel(t *testing.T) {
	ll, err := logrus.ParseLevel(logrusLevel)
	require.NoError(t, err)
	logrus.SetLevel(ll)
	logrus.Infof("logrus level has been set to %s", logrus.GetLevel().String())
}

// CreateTempFile creates a temp file
func CreateTempFile(t *testing.T) *os.File {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "prefix-")
	require.NoError(t, err)
	return tmpFile
}
