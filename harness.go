package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime/trace"

	"github.com/chromedp/chromedp"
)

var traceDir string = "traces"

func main() {
	// defer trace.Stop()
	// trace.Start(os.Stderr)
	fmt.Println("Hello World!")
	createDirIfNotExist(traceDir)
	runOne()
	chromeDPscreenshot()
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

func chromeDPscreenshot() {
	// create context
	// ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithDebugf(log.Printf))
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// capture screenshot of an element
	var buf []byte
	if err := chromedp.Run(ctx, traceScreenshot(`http://127.0.0.1:8080/trace`, &buf)); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile("elementScreenshot.png", buf, 0644); err != nil {
		log.Fatal(err)
	}

}

func traceScreenshot(urlstr string, res *[]byte) chromedp.Tasks {
	// chromedp.Sleep(15 * time.Second),
	loadingSelector := "body > overlay"
	traceViewerId := "#trace-viewer"
	mouseModeSelector := "#track_view_container > tr-ui-timeline-track-view > tr-ui-b-mouse-mode-selector"
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.WaitNotPresent(loadingSelector, chromedp.ByQuery),
		chromedp.WaitVisible(traceViewerId, chromedp.ByID),
		chromedp.WaitVisible(mouseModeSelector, chromedp.ByQuery),
		chromedp.SetAttributeValue(mouseModeSelector, "style.display", "none"),
		chromedp.Screenshot(traceViewerId, res, chromedp.NodeVisible, chromedp.ByID),
	}
}
