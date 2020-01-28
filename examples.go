package main

import (
	"fmt"
	"reflect"
	"time"
)

// GatherExamples is the epitome of a hack. I wanted to learn reflection, so
// this returns the Method for all exported functions attached to the Examples struct
func GatherExamples() ([]reflect.Value, []string) {
	exVal := reflect.ValueOf(Examples{})
	exType := reflect.TypeOf(Examples{})
	mVals := make([]reflect.Value, exType.NumMethod())
	mNames := make([]string, exType.NumMethod())
	for i := 0; i < exType.NumMethod(); i++ {
		mVals[i] = exVal.Method(i)
		mNames[i] = exType.Method(i).Name

	}

	return mVals, mNames
}

type Examples struct {
}

func (e Examples) Hello() {
	fmt.Println("Hello World!")
}

func (e Examples) BrokenHello() {
	fmt.Println("broken World!")
	go func() {
		time.Sleep(1)
		fmt.Println("Hello World!")
	}()
}

func (e Examples) FixedHello() {
	fmt.Println("fixed World!")
	go func() {
		time.Sleep(1)
		fmt.Println("Hello World!")
	}()
	time.Sleep(2)
}
