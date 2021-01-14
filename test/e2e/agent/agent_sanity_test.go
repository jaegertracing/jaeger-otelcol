// +build agent_smoke

package e2e

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/jaegertracing/jaeger-otelcol/test/e2e"
	"github.com/jaegertracing/jaeger-otelcol/test/tools/tracegen"
)

type AgentSanityTestSuite struct {
	suite.Suite
}

var t *testing.T
var logger zap.SugaredLogger

func (suite *AgentSanityTestSuite) SetupSuite() {
	logger = e2e.GetLogger(suite.T())
}

func (suite *AgentSanityTestSuite) TearDownSuite() {
	logger.Infof("In teardown suite")
}

func TestAgentSanityTestSuite(t *testing.T) {
	suite.Run(t, new(AgentSanityTestSuite))
}

func (suite *AgentSanityTestSuite) BeforeTest(suiteName, testName string) {
	t = suite.T()
	logger.Debugf("In Before for %s", suite.T().Name())
}

func (suite *AgentSanityTestSuite) AfterTest(suiteName, testName string) {
	logger.Debug("In AfterTest for %s", suite.T().Name())
}

func (suite *AgentSanityTestSuite) TestAgentSanity() {
	// Start the agent
	agentExecutable := "../../../builds/agent/jaeger-otel-agent"
	agentConfigFileName := "./config/jaeger-agent-config.yaml"
	metricsPort := e2e.GetFreePort(t)
	logger.Infof("Using metrics port %s", metricsPort)

	loggerOutputFile := e2e.CreateTempFile(t)
	logger.Infof("Using log file %s", loggerOutputFile.Name())
	agent := e2e.StartCollector(t, logger, agentExecutable, agentConfigFileName, loggerOutputFile, metricsPort)
	defer agent.Process.Kill()

	// Create some traces. Each trace created by tracegen will have 2 spans
	traceCount := 5
	expectedSpanCount := 2 * traceCount
	serviceName := "agent-sanity-test" + strconv.Itoa(time.Now().Nanosecond())
	tracegen.CreateJaegerTraces(t, 1, traceCount, 0, serviceName)

	// Check the metrics to verify that the agent received and then sent the number of spans expected
	metricsEndpoint := "http://localhost:" + metricsPort + "/metrics"
	receivedSpansCounter := e2e.GetPrometheusCounter(t, metricsEndpoint, "otelcol_receiver_accepted_spans")
	sentSpansCounter := e2e.GetPrometheusCounter(t, metricsEndpoint, "otelcol_exporter_sent_spans")
	require.Equal(t, expectedSpanCount, int(receivedSpansCounter))
	require.Equal(t, expectedSpanCount, int(sentSpansCounter))

	// Don't do a defer, only remove the log if the test passes
	os.Remove(loggerOutputFile.Name())
}
