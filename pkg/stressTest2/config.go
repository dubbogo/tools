package stressTest

import (
	"time"
)

type StressTestConfig struct {
	tps                     int
	duration                time.Duration
	successRateGaugeHandler func(rate float64)
	errorCounterHandler     func(errCount uint32)
}

func (c *StressTestConfig) Start(f func() error) {
	newTester(c, f).Test()
}

type StressTestConfigBuilder struct {
	s *StressTestConfig
}

func (s *StressTestConfigBuilder) SetDuration(duration string) *StressTestConfigBuilder {
	var err error
	s.s.duration, err = time.ParseDuration(duration)
	if err != nil {
		panic(err)
	}
	return s
}

func (s *StressTestConfigBuilder) SetTPS(tps int) *StressTestConfigBuilder {
	s.s.tps = tps
	return s
}

func (s *StressTestConfigBuilder) SetSuccessRateGaugeHandler(successRateGaugeHandler func(rate float64)) *StressTestConfigBuilder {
	s.s.successRateGaugeHandler = successRateGaugeHandler
	return s
}

func (s *StressTestConfigBuilder) SetErrorCounterHandler(errorCounterHandler func(rate uint32)) *StressTestConfigBuilder {
	s.s.errorCounterHandler = errorCounterHandler
	return s
}

func (s *StressTestConfigBuilder) Build() *StressTestConfig {
	return s.s
}

func NewStressTestConfigBuilder() *StressTestConfigBuilder {
	return &StressTestConfigBuilder{
		s: &StressTestConfig{
			tps: 0,
		},
	}
}
