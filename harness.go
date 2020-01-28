package main

import (
	"fmt"
	gwChrome "github.com/sensepost/gowitness/chrome"
	"net/url"
	"os"
	"runtime/trace"
)

var traceDir string = "traces"

func main() {
	// defer trace.Stop()
	// trace.Start(os.Stderr)
	fmt.Println("Hello World!")
	createDirIfNotExist(traceDir)
	runOne()
	screenshotTrace()
}

func createDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func runOne() {
	f, err := os.Create(traceDir + "/out.trace")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = trace.Start(f)
	if err != nil {
		panic(err)
	}

	Hello()

	trace.Stop()
}

func screenshotTrace() {

	chrome := &gwChrome.Chrome{
		Resolution:    `800x600`,
		ChromeTimeout: 30,
	}
	chrome.Setup()

	u, err := url.ParseRequestURI(`http://127.0.0.1:8080/trace`)
	if err != nil {
		panic(err)
	}
	chrome.ScreenshotURL(u, `trace.png`)
}
