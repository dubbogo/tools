package stressTest

import (
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"fmt"
	"sync"
	"time"
)

type StressTester struct {
	config   *StressTestConfig
	testFunc func() error
}

func (s *StressTester) Test() {
	counter := 0
	errorCounter := 0
	rtList := make([]int64, 0, s.config.tps)
	//var rtList []int64
	lock := sync.Mutex{}

	// controller
	startTestTime := time.Now()
	// send first period
	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C
		for i := 0; i < s.config.tps; i++ {
			if i == 0 {
				go func() {
					startTime := time.Now()
					err := s.testFunc()
					lock.Lock()
					rtList = append(rtList, int64(time.Now().Sub(startTime)))
					if err != nil {
						errorCounter++
					} else {
						counter++
					}
					lock.Unlock()
				}()
				continue
			}
			go func() {
				err := s.testFunc()
				lock.Lock()
				if err != nil {
					logger.Error(err)
					errorCounter++
				} else {
					counter++
				}
				lock.Unlock()
			}()
		}
		// calculate result rt
		lock.Lock()
		avgRT := int64(9999999999)
		if len(rtList) != 0 {
			avgRT = average(rtList)
		}
		// read
		tempCounter := counter
		tempErrorCounter := errorCounter

		// set to zero
		errorCounter = 0
		counter = 0
		rtList = []int64{}
		lock.Unlock()
		fmt.Println("average rt = ", avgRT, " tps = ", tempCounter, "errorCount = ", tempErrorCounter)
		if tempCounter+tempErrorCounter != 0 {
			successRate := (float64(tempCounter) / float64(tempCounter+tempErrorCounter))
			s.config.successRateGaugeHandler(successRate)
			fmt.Printf(" success rate = %f\n", successRate)
			if tempErrorCounter != 0 {
				s.config.errorCounterHandler(uint32(tempErrorCounter))
			}
		}

		if time.Now().Sub(startTestTime) >= s.config.duration {
			// close all
			return
		}
	}

}

func newTester(config *StressTestConfig, testFunc func() error) *StressTester {
	return &StressTester{config: config, testFunc: testFunc}
}

func average(xs []int64) (avg int64) {
	sum := int64(0)
	switch len(xs) {
	case 0:
		avg = 0
	default:
		for _, v := range xs {
			sum += v
		}
		avg = sum / int64(len(xs))
	}
	return
}
