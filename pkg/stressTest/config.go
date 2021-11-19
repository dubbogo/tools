package stressTest

import (
	"time"
)

type StressTestConfig struct {
	parallel int
	tps      int
	duration time.Duration
}

func (c *StressTestConfig) Start(f func()) {
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

func (s *StressTestConfigBuilder) SetParallel(num int) *StressTestConfigBuilder {
	s.s.parallel = num
	return s
}

func (s *StressTestConfigBuilder) Build() *StressTestConfig {
	return s.s
}

func NewStressTestConfigBuilder() *StressTestConfigBuilder {
	return &StressTestConfigBuilder{
		s: &StressTestConfig{
			parallel: 1,
			tps:      0,
		},
	}
}
