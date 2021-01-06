package e2e

import (
	"crypto/tls"
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

func GetPrometheusCounter(t *testing.T, metricsEndpoint, metricName string) float64 {
	counter := GetPrometheusMetric(t, metricsEndpoint, metricName)
	return *counter.Metric[0].Counter.Value
}

func GetPrometheusMetric(t *testing.T, metricsEndpoint, metricName string) pcm.MetricFamily {
	allMetrics := GetPrometheusMetrics(t, metricsEndpoint)
	return allMetrics[metricName]
}


// This code is mostly copied from https://github.com/prometheus/prom2json except it
// returns MetricFamily objects as that is more useful than JSON for tests.
func GetPrometheusMetrics(t *testing.T, metricsEndpoint string) map[string]pcm.MetricFamily {
	mfChan := make(chan *pcm.MetricFamily, 1024)
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	
	err := prom2json.FetchMetricFamilies(metricsEndpoint, mfChan, transport)
	require.NoError(t, err)
	result := map[string]pcm.MetricFamily{}
	for mf := range mfChan {
		result[*mf.Name] = *mf
	}

	return result
}

func SetLogrusLevel(t *testing.T) {
	ll, err := logrus.ParseLevel(logrusLevel)
	require.NoError(t, err)
	logrus.SetLevel(ll)
	logrus.Infof("logrus level has been set to %s", logrus.GetLevel().String())
}

func CreateTempFile(t *testing.T) *os.File {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "prefix-")
	require.NoError(t, err)
	return tmpFile
}
