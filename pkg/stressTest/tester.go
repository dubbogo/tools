package stressTest

import (
	"fmt"
	"sync"
	"time"
)

type StressTester struct {
	config   *StressTestConfig
	testFunc func()
}

func (s *StressTester) Test() {
	counter := 0
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
						s.testFunc()
						rtList = append(rtList, int64(time.Now().Sub(startTime)))

						lock.Lock()
						counter++
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
					s.testFunc()

					lock.Lock()
					counter++
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
		counter = 0
		lock.Unlock()
		fmt.Println("average rt = ", avgRT, " tps = ", tempCounter)

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

func newTester(config *StressTestConfig, testFunc func()) *StressTester {
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
