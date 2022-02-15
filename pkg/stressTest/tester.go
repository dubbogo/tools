package stressTest

import (
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
	taskChan := make(chan bool, s.config.tps)
	closeChan := make(chan struct{})
	rtList := make([]int64, 0, s.config.tps)
	allRTList := make([]int64, 0, s.config.duration/time.Second)
	//var rtList []int64
	lock := sync.Mutex{}

	for i := 0; i < s.config.parallel; i++ {
		if i == 0 {
			go func() {
				for {
					select {
					case <-closeChan:
						return
					case <-taskChan:

						startTime := time.Now()
						err := s.testFunc()
						rtList = append(rtList, int64(time.Now().Sub(startTime)))

						lock.Lock()
						if err != nil {
							errorCounter++
						} else {
							counter++
						}
						lock.Unlock()
					}

				}
			}()
			continue
		}
		go func() {
			for {
				select {
				case <-closeChan:
					return
				case <-taskChan:
					err := s.testFunc()

					lock.Lock()
					if err != nil {
						errorCounter++
					} else {
						counter++
					}
					lock.Unlock()
				}
			}
		}()
	}

	// controller
	startTestTime := time.Now()
	// send first period
	for i := 0; i < s.config.tps; i++ {
		taskChan <- true
	}

	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C

		// empty the taskChan
	LOOP:
		for {
			select {
			case <-taskChan:
			default:
				break LOOP
			}
		}

		// calculate result rt
		lock.Lock()
		avgRT := int64(999999)
		if len(rtList) != 0 {
			avgRT = average(rtList)
			allRTList = append(allRTList, avgRT)
			rtList = make([]int64, 0, s.config.tps)
		}
		tempCounter := counter
		tempErrorCounter := errorCounter
		counter = 0
		errorCounter = 0
		lock.Unlock()
		fmt.Println("average rt = ", avgRT, " tps = ", tempCounter)
		if tempCounter+tempErrorCounter != 0 {
			fmt.Printf(" success rate = %f\n", (float32(tempCounter) / float32(tempCounter+tempErrorCounter)))
		}

		if time.Now().Sub(startTestTime) >= s.config.duration {
			// close all
			close(closeChan)
			return
		}

		// send next period
		for i := 0; i < s.config.tps; i++ {
			taskChan <- true
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
