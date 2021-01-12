// +build agent_smoke

package e2e

import (
	"fmt"

	"io/ioutil"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jaegertracing/jaeger-otelcol/test/e2e"
	"github.com/jaegertracing/jaeger-otelcol/test/tools/tracegen"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AgentSanityTestSuite struct {
	suite.Suite
}

var t *testing.T

func (suite *AgentSanityTestSuite) SetupSuite() {
	e2e.SetLogrusLevel(suite.T())
}

func (suite *AgentSanityTestSuite) TearDownSuite() {
	logrus.Infof("In teardown suite")
}

func TestAgentSanityTestSuite(t *testing.T) {
	suite.Run(t, new(AgentSanityTestSuite))
}

func (suite *AgentSanityTestSuite) BeforeTest(suiteName, testName string) {
	t = suite.T()
	logrus.Infof("In Before for %s", suite.T().Name())
}

func (suite *AgentSanityTestSuite) AfterTest(suiteName, testName string) {
	logrus.Infof("In AfterTest for %s", suite.T().Name())
}

func (suite *AgentSanityTestSuite) TestAgentSanity() {
	// Start the agent
	agentExecutable := "../../../builds/agent/jaeger-otel-agent"
	agentConfigFileName := "./config/jaeger-agent-config.yaml"
	metricsPort := "8888"

	loggerOutputFile := e2e.CreateTempFile(t)
	defer os.Remove(loggerOutputFile.Name())
	agent := e2e.StartCollector(t, agentExecutable, agentConfigFileName, loggerOutputFile, metricsPort)
	defer agent.Process.Kill()

	// Create some traces. Each trace created by tracegen will have 2 spans
	traceCount := 5
	expectedSpanCount := 2 * traceCount
	serviceName := "agent-sanity-test" + strconv.Itoa(time.Now().Nanosecond())
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
