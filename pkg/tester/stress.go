package tester

import (
	"context"
	"fmt"
	"math"
	"time"
)

import (
	"go.uber.org/atomic"
)

const (
	Inf = math.MaxInt32
)

type TestFnType func(userId int)
type TestFnWithErrorType func(userId int) error

type StressTester struct {
	// config
	verbose           bool
	duration          time.Duration
	tps               int
	testFn            TestFnType
	testFnWithError   TestFnWithErrorType
	enableSuccessRate bool

	startedTime    *time.Time
	endedTime      *time.Time
	transactionNum *atomic.Int32
	rt             *atomic.Float64
	successRate    *atomic.Float64
}

func NewStressTester() *StressTester {
	return &StressTester{}
}

func (t *StressTester) Run() {
	defer func() {
		now := time.Now()
		t.endedTime = &now
		t.log("The stress tester finished, %d requests are sent", t.GetTransactionNum())
	}()

	if t.startedTime != nil || t.transactionNum != nil || t.endedTime != nil {
		fmt.Println("Run stress tester after reset it.")
		return
	}

	now := time.Now()
	t.startedTime = &now
	t.transactionNum = new(atomic.Int32)
	t.rt = new(atomic.Float64)

	ctx, cancel := context.WithTimeout(context.Background(), t.duration)
	defer cancel()

	t.user(ctx)
}

func (t *StressTester) user(ctx context.Context) {
	counter := t.tps
	timer := time.After(1 * time.Second)

	remainingRequestNum := new(atomic.Int32)

	for {
		select {
		case <-ctx.Done():
			// run out of the time, but it should do graceful shutdown
			t.log("Run out of the time, stress tester is going to shutdown...")
			for remainingRequestNum.Load() > 0 {
				t.verboseLog("Waiting for all tests to be finished, recheck after 500ms...")
				time.Sleep(500 * time.Millisecond)
			}
			return
		case <-timer:
			// trigger once in one second to detect the TPS limitation
			t.log("[Warning] The TPS is too big to send all requests in one second.")
			counter = t.tps
			timer = time.After(1 * time.Second)
		default:
		}

		go func() {
			startedTime := time.Now()
			defer func() {
				t.transactionNum.Add(1)
				t.rt.Add(time.Now().Sub(startedTime).Seconds())
				remainingRequestNum.Add(-1)
			}()

			remainingRequestNum.Add(1)
			if t.enableSuccessRate {
				err := t.testFnWithError(0)
				if err != nil {

				}
			} else {
				t.testFn(0)
			}
		}()

		counter--
		if counter == 0 {
			select {
			case <-timer:
				counter = t.tps
				timer = time.After(1 * time.Second)
			}
		}
	}
}

func (t *StressTester) GetTPS() float64 {
	if t.startedTime == nil || t.transactionNum == nil {
		return -1
	}
	return float64(t.transactionNum.Load()) / t.GetElapsedTimeSeconds()
}

func (t *StressTester) GetElapsedTimeSeconds() float64 {
	if t.startedTime == nil {
		return -1
	}
	if t.endedTime == nil {
		return time.Now().Sub(*t.startedTime).Seconds()
	}
	return t.endedTime.Sub(*t.startedTime).Seconds()
}

func (t *StressTester) GetSuccessRate() float64 {
	return t.successRate.Load()
}

func (t *StressTester) GetAverageRTSeconds() float64 {
	if t.rt == nil || t.transactionNum == nil {
		return -1
	}
	return t.rt.Load() / float64(t.transactionNum.Load())
}

func (t *StressTester) GetTransactionNum() int32 {
	if t.transactionNum == nil {
		return -1
	}
	return t.transactionNum.Load()
}

func (t *StressTester) SetVerbose(v bool) *StressTester {
	t.verbose = v
	return t
}

// SetUserNum is depreciated for now, please use TPS instead, it will be removed from older version.
// TODO(justxuewei): compatible with previous version
func (t *StressTester) SetUserNum(_ int) *StressTester {
	t.log("please note that user num is not supported for now, please use TPS instead")
	return t
}

func (t *StressTester) SetDuration(d string) *StressTester {
	duration, err := time.ParseDuration(d)
	if err != nil {
		panic(nil)
	}
	t.duration = duration
	return t
}

func (t *StressTester) SetTPS(tps int) *StressTester {
	t.tps = tps
	return t
}

func (t *StressTester) SetTestFn(fn TestFnType) *StressTester {
	t.testFn = fn
	return t
}

func (t *StressTester) EnableSuccessRate() *StressTester {
	t.enableSuccessRate = true
	return t
}

func (t *StressTester) SetTestFnWithError(fn TestFnWithErrorType) *StressTester {
	t.testFnWithError = fn
	return t
}

func (t *StressTester) Reset() {
	t.startedTime = nil
	t.endedTime = nil
	t.transactionNum = nil
	t.rt = nil
}

func (t *StressTester) verboseLog(msg string, args ...interface{}) {
	if t.verbose {
		t.log(msg, args...)
	}
}

func (t *StressTester) log(msg string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(msg, args...))
}
