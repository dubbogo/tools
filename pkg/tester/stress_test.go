package tester

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestStressTester(t *testing.T) {

	t.Run("tps", func(t *testing.T) {

		testFn := func(userId int) {
			fmt.Printf("testFn started, userId = %d\n", userId)
			time.Sleep(20 * time.Second)
			fmt.Printf("testFn finished, userId = %d\n", userId)
		}

		tester := NewStressTester()
		tester.
			SetVerbose(true).
			SetTPS(10).
			SetDuration("10s").
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
			SetDuration("10s").
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
