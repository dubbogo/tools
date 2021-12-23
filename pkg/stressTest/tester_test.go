package stressTest

import (
	"testing"
	"time"
)

func TestTester(t *testing.T) {
	config := new(StressTestConfig)
	config.tps = 100
	config.parallel = 5
	config.duration = 1 * time.Minute

	tester := newTester(config, func() {
		time.Sleep(10 * time.Second)
	})

	tester.Test()
}
