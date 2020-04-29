package main

import (
	"fmt"
	"time"

	"github.com/dty1er/traffic-shaping-io/tokenbucket"
)

func main() {
	b := tokenbucket.New(10, 10)
	fmt.Println("start")
	time.Sleep(2 * time.Second)
	b.Take(50)
}
