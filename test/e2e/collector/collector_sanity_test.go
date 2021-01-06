// +build collector_smoke

package e2e

import (
	"fmt"
	"github.com/jaegertracing/jaeger-otelcol/test/tools/tracegen"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jaegertracing/jaeger-otelcol/test/e2e"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CollectorSanityTestSuite struct {
	suite.Suite
}

var t *testing.T

func (suite *CollectorSanityTestSuite) SetupSuite() {
	e2e.SetLogrusLevel(suite.T())
}

func (suite *CollectorSanityTestSuite) TearDownSuite() {
	logrus.Infof("In teardown suite")
}

func TestCollectorSanityTestSuite(t *testing.T) {
	suite.Run(t, new(CollectorSanityTestSuite))
}

func (suite *CollectorSanityTestSuite) BeforeTest(suiteName, testName string) {
	t = suite.T()
	logrus.Infof("In Before for %s", suite.T().Name())
}

func (suite *CollectorSanityTestSuite) AfterTest(suiteName, testName string) {
	logrus.Infof("In AfterTest for %s", suite.T().Name())
}

// TODO can we combine this and the TestAgentSanity test?  What is different besides executable, configfile, servicename, traceCount
func (suite *CollectorSanityTestSuite) TestCollectorSanity() {
	// Start the collector
	colletorExecutable := "../../../builds/collector/jaeger-otel-collector"
	collectorConfigFileName := "./config/jaeger-collector-config.yaml"
	metricsPort := e2e.GetFreePort(t)
	logrus.Infof("Using metrics port %s", metricsPort)

	loggerOutputFile := e2e.CreateTempFile(t)
	defer os.Remove(loggerOutputFile.Name())
	agent := e2e.StartCollector(t, colletorExecutable, collectorConfigFileName, loggerOutputFile, metricsPort)
	defer agent.Process.Kill()

	// Create some traces. Each trace created by tracegen will have 2 spans
	traceCount := 5
	expectedSpanCount := 2 * traceCount
	serviceName := "collector-sanity-test" + strconv.Itoa(time.Now().Nanosecond())
	tracegen.CreateJaegerTraces(t, 1, traceCount, 0, serviceName)

	// This could be changed to logrus.Debugf if we can stop logrus from eating newlines
	if logrus.GetLevel() == logrus.DebugLevel {
		log, err := ioutil.ReadFile(loggerOutputFile.Name())
		require.NoError(t, err)
		fmt.Printf("%s", log)
	}

	// Check the metrics to verify that the agent received and then sent the number of spans expected
	metricsEndpoint := "http://localhost:" + metricsPort + "/metrics"
	receivedSpansCounter := e2e.GetPrometheusCounter(t, metricsEndpoint, "otelcol_receiver_accepted_spans")
	sentSpansCounter := e2e.GetPrometheusCounter(t, metricsEndpoint, "otelcol_exporter_sent_spans")
	require.Equal(t, expectedSpanCount, int(receivedSpansCounter))
	require.Equal(t, expectedSpanCount, int(sentSpansCounter))
}
