package main

import (
	"fmt"
	"log"
	"time"
)

var startTime = time.Now()

func main() {
	bucket()
}

func bucket() {
	quantum := int64(1)
	rate := float64(2)
	capacity := int64(2)
	fillInterval := time.Duration(1e9 * float64(quantum) / rate)
	fmt.Println(fillInterval)
	count := int64(1) //每次从桶中取token的数量
	availableTokens := int64(2)
	var latestTick int64
	t := time.NewTicker(400 * time.Millisecond)
	for e := range t.C {
		tick := int64(time.Now().Sub(startTime) / fillInterval)
		if availableTokens >= capacity {
			fmt.Printf("ok use [%d]", availableTokens)
		} else {
			availableTokens += (tick - latestTick) * quantum
			if availableTokens > capacity {
				availableTokens = capacity
			}
			latestTick = tick
		}
		if availableTokens <= 0 {
			log.Println("no availableTokens")
		}else{
			if count > availableTokens{
				count = availableTokens
			}
			availableTokens -= count
			fmt.Printf("[%+v] === > availableTokens has [%d]",e,availableTokens)
			fmt.Println()
		}
	}
}
