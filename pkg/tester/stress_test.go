package tester

import (
	"testing"
	"time"
)

func TestStressTester(t *testing.T) {

	t.Run("tps", func(t *testing.T) {

		testFn := func(userId int) {
			time.Sleep(20 * time.Second)
		}

		tester := NewStressTester()
		tester.
			SetVerbose(true).
			SetTPS(10).
			SetDuration(10 * time.Second).
			SetUserNum(2).
			SetTestFn(testFn).
			Run()

		t.Log("TPS:", tester.GetTPS())
		t.Log("Elapsed Time:", tester.GetElapsedTimeSeconds())

		tester.Reset()
		t.Log("TPS:", tester.GetTPS())
		t.Log("Elapsed Time:", tester.GetElapsedTimeSeconds())
	})

}
