package e2e

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

var (
	logrusLevel = getStringEnv("LOGRUS_LEVEL", "info")
)

// Metric can contain the value of a prometheus metric
type Metric struct {
	Key      string
	JSONPart string
	Value    string
}

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

// GetMetric returns the value of the named metric
func GetMetric(t *testing.T, metricsEndpoint string, metricKey string) Metric {
	for _, m := range GetMetrics(t, metricsEndpoint) {
		if m.Key == strings.TrimSpace(metricKey) {
			return m
		}
	}

	logrus.Warnf("Could not find metric %s at endpoint %s", metricKey, metricsEndpoint)
	emptyMetric := Metric{}
	return emptyMetric
}

// GetMetrics returns all metrics from the endpoint
func GetMetrics(t *testing.T, metricsEndpoint string) []Metric {
	httpClient := http.Client{Timeout: 5 * time.Second}

	request, err := http.NewRequest(http.MethodGet, metricsEndpoint, nil)
	require.NoError(t, err)
	response, err := httpClient.Do(request)
	require.NoError(t, err)

	body, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)

	lines := strings.Split(string(body), "\n")
	metrics := []Metric{}
	for _, line := range lines {
		if !strings.HasPrefix(line, "#") && strings.TrimSpace(line) != "" {
			openBracket := strings.Index(line, "{")
			closeBracket := strings.Index(line, "}") + 1
			key := line[:openBracket]
			jsonPart := line[openBracket:closeBracket]
			value := line[closeBracket:]

			metric := Metric{
				Key:      key,
				JSONPart: jsonPart,
				Value:    strings.TrimSpace(value),
			}
			metrics = append(metrics, metric)
		}
	}

	return metrics
}

// SetLogrusLevel can be used to set the logging level
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
