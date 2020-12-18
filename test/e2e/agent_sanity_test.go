// +build agent_smoke

package e2e

import (
	"bytes"
	"fmt"
	"index/suffixarray"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/jaegertracing/jaeger-otelcol/test/tools/tracegen"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AgentSanityTestSuite struct {
	suite.Suite
}

func (suite *AgentSanityTestSuite) SetupSuite() {
	setLogrusLevel(suite.T())
}

func (suite *AgentSanityTestSuite) TearDownSuite() {
	logrus.Infof("In teardown wuite")
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
	agentExecutable := "../../builds/agent/jaeger-otel-agent"
	agentConfigFileName := "./config/jaeger-agent-config.yaml"
	metricsPort := "8888"

	var agentLoggerOutput bytes.Buffer
	agent := StartCollector(t, agentExecutable, agentConfigFileName, &agentLoggerOutput, metricsPort)
	defer agent.Process.Kill()

	// Create some traces. Each trace created by tracegen will have 2 spans
	traceCount := 5
	expectedSpanCount := 2 * traceCount
	serviceName := "agent-sanity-test" + strconv.Itoa(time.Now().Nanosecond())
	tracegen.CreateJaegerTraces(t, 1, traceCount, 0, serviceName)

	// This could be changed to logrus.Debugf if we can stop logrus from eating newlines
	if logrus.GetLevel() == logrus.DebugLevel {
		fmt.Printf("%s", agentLoggerOutput.String())
	}

	// Check the agent output for service name and Span entries
	require.Contains(t, agentLoggerOutput.String(), "service.name: STRING("+serviceName+")")

	spanExpression := regexp.MustCompile("Span #")
	index := suffixarray.New(agentLoggerOutput.Bytes())
	results := index.FindAllIndex(spanExpression, -1)
	require.Equal(t, expectedSpanCount, len(results))

	// Check the metrics to verify that the agent received and then sent the number of spans expected
	metricsEndpoint := "http://localhost:" + metricsPort + "/metrics"
	receivedSpansMetric := getMetric(t, metricsEndpoint, "otelcol_receiver_accepted_spans")
	sentSpansMetric := getMetric(t, metricsEndpoint, "otelcol_exporter_sent_spans")
	require.Equal(t, strconv.Itoa(expectedSpanCount), receivedSpansMetric.Value)
	require.Equal(t, strconv.Itoa(expectedSpanCount), sentSpansMetric.Value)
}
