package main

import (
	"fmt"
)

func InitiateGoroutine() {
	fmt.Println("Hello World!")

	for i := 0; i < 10; i++ {
		go fmt.Println(i)
	}
}
