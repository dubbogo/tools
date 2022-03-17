package tester

import (
	"errors"
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
			SetDuration("10s").
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

func TestStressTesterWithError(t *testing.T) {

	t.Run("tpsWithErr", func(t *testing.T) {

		testFnWithError := func(userId int) error {
			if userId%2 == 0 {
				return errors.New("err")
			}
			return nil
		}

		tester := NewStressTester()
		tester.
			SetVerbose(true).
			SetTPS(10).
			SetDuration(10 * time.Second).
			SetUserNum(2).
			SetTestFnWithError(testFnWithError).
			EnableSuccessRate().
			Run()

		t.Log("TPS:", tester.GetTPS())
		t.Log("Elapsed Time:", tester.GetElapsedTimeSeconds())
		t.Log("Success Rate:", tester.GetSuccessRate())

		tester.Reset()
		t.Log("TPS:", tester.GetTPS())
		t.Log("Elapsed Time:", tester.GetElapsedTimeSeconds())
	})

}
