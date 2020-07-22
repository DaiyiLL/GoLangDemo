package main

import (
	"fmt"
	"runtime"
	"time"
)

func aaa()  {
	for {
		time.Sleep(time.Millisecond * 200)
		fmt.Println("-----------------------\n")
	}
}

func main00()  {
	go func() {
		fmt.Println("----------------------1\n")
		go aaa()

		fmt.Println("----------------------2\n")
		return
	}()

	for {
		runtime.GC()
	}
}
