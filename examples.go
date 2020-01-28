package main

import (
	"fmt"
	"time"
)

func Hello() {
	fmt.Println("Hello World!")
}

func BrokenHello() {
	go func() {
		time.Sleep(1)
		fmt.Println("Hello World!")
	}()
}

func FixedHello() {
	go func() {
		time.Sleep(1)
		fmt.Println("Hello World!")
	}()
	time.Sleep(2)
}
