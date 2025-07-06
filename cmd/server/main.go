package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()

	time.Sleep(1 * time.Second)

	elapsed := time.Since(now).Seconds()

	fmt.Println(elapsed)
}
