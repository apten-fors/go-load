package main

// #include <time.h>
import "C"

import (
	"flag"
	"runtime"
	"sync"
	"time"
)

var startTime = time.Now()
var startTicks = C.clock()
var WaitFactor = 1

func CpuUsagePercent(samplingRate float64) float64 {
	clockSeconds := float64(C.clock()-startTicks) / float64(C.CLOCKS_PER_SEC)
	realSeconds := time.Since(startTime).Seconds()
	if samplingRate > 0 && realSeconds >= samplingRate {
		startTime = time.Now()
		startTicks = C.clock()
	}
	return clockSeconds / realSeconds * 100
}

func main() {
	defer func() {
		if errP := recover(); errP != nil {
			return
		}
	}()
	tcpu := flag.Float64("cpu", 30.0, "targetCPU")
	rate := flag.Float64("rate", 0.0, "samplingRate")
	flag.Parse()

	var wg sync.WaitGroup
	defer func() {
		if errP := recover(); errP != nil {
		}
	}()
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		targetCPU := *tcpu
		samplingRate := *rate
		for i := 0; i < runtime.NumCPU(); i++ {
			go func() {
				for {
					cpuLoad := CpuUsagePercent(samplingRate)
					if cpuLoad >= targetCPU {
						waitTime := time.Second / time.Duration(WaitFactor)
						time.Sleep(waitTime)
					}
				}
				defer wg.Done()
			}()
		}
	}(&wg)
	wg.Wait()
}
