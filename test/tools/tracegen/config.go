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

package tracegen

import (
	"fmt"
	"testing"
	"sync"
	"sync/atomic"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/require"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	jaegerZap "github.com/uber/jaeger-client-go/log/zap"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	"go.uber.org/zap"
)

var logger, _ = zap.NewDevelopment()

// Config describes the test scenario.
type Config struct {
	Workers  int
	Traces   int
	Marshal  bool
	Debug    bool
	Firehose bool
	Pause    time.Duration
	Duration time.Duration
	Service  string
}

func CreateJaegerTraces(t *testing.T, workers, traces int, duration time.Duration, serviceName string) {
	config := &Config{
		Workers:  workers,
		Traces:   traces,
		Marshal:  false,
		Debug:    false,
		Firehose: false,
		Pause:    0,
		Duration: duration,
		Service:  serviceName,
	}

	metricsFactory := prometheus.New()
	traceCfg := &jaegerConfig.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		RPCMetrics: true,
	}
	traceCfg, err := traceCfg.FromEnv()

	require.NoError(t, err)

	tracer, tCloser, err := traceCfg.NewTracer(
		jaegerConfig.Metrics(metricsFactory),
		jaegerConfig.Logger(jaegerZap.NewLogger(logger)),
	)
	require.NoError(t, err)
	defer tCloser.Close()

	opentracing.InitGlobalTracer(tracer)
	logger.Info("Initialized global tracer")

	Run(config, logger)

	logger.Info("Waiting 1.5sec for metrics to flush")
	time.Sleep(3 * time.Second / 2)
}

// Run executes the test scenario.
func Run(c *Config, logger *zap.Logger) error {
	if c.Duration > 0 {
		c.Traces = 0
	} else if c.Traces <= 0 {
		return fmt.Errorf("either `traces` or `duration` must be greater than 0")
	}

	wg := sync.WaitGroup{}
	var running uint32 = 1
	for i := 0; i < c.Workers; i++ {
		wg.Add(1)
		w := worker{
			id:       i,
			traces:   c.Traces,
			marshal:  c.Marshal,
			debug:    c.Debug,
			firehose: c.Firehose,
			pause:    c.Pause,
			duration: c.Duration,
			running:  &running,
			wg:       &wg,
			logger:   logger.With(zap.Int("worker", i)),
		}

		go w.simulateTraces()
	}
	if c.Duration > 0 {
		time.Sleep(c.Duration)
		atomic.StoreUint32(&running, 0)
	}
	wg.Wait()
	return nil
}
