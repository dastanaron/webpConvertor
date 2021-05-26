package main

import (
	"fmt"
	"time"
)

func main() {
	startTime := time.Now()

	endTime := time.Now()

	fmt.Println("Time:", endTime.Sub(startTime))
}
