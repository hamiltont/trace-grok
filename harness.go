package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"runtime/trace"
	"syscall"
	"time"

	"github.com/chromedp/chromedp"
)

var traceDir string = "traces"

func main() {
	// defer trace.Stop()
	// trace.Start(os.Stderr)
	createDirIfNotExist(traceDir)

	methods, names := GatherExamples()

	for i := 0; i < len(names); i++ {
		runOne(names[i], methods[i])
	}

	fmt.Printf("Done running one, %d goroutines remaining\n", runtime.NumGoroutine())
	// pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
	// time.Sleep(5 * time.Second)
	// fmt.Printf("Done sleeping, %d goroutines remaining\n", runtime.NumGoroutine())

	// pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
	// pprof.Lookup("threadcreate").WriteTo(os.Stdout, 1)
}

func createDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func runOne(name string, example reflect.Value) {
	fmt.Printf("Tracing %s (%d)\n", name, runtime.NumGoroutine())

	tracePath := traceDir + "/" + name + ".trace"
	traceImgPath := traceDir + "/" + name + ".png"

	f, err := os.Create(tracePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = trace.Start(f)
	if err != nil {
		panic(err)
	}

	example.Call(nil)
	// (Examples{}).Hello()
	trace.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// cmd := exec.CommandContext(ctx, "ls")
	// TODO this appears to be the broken line, casuses vim-go to hang forever
	fmt.Printf("Before running trace http process (%d)\n", runtime.NumGoroutine())
	cmd := exec.CommandContext(ctx, "go", "tool", "trace", "-http=127.0.0.1:32000", tracePath)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Env = append(os.Environ(), "BROWSER=true")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Started trace http process (%d)\n", runtime.NumGoroutine())

	fmt.Printf("Taking screenshot (%d)\n", runtime.NumGoroutine())
	chromeDPscreenshot(traceImgPath)

	fmt.Println("Killing trace")

	syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)

	cmd.Process.Signal(os.Interrupt)
	if err := cmd.Process.Kill(); err != nil {
		log.Fatal("failed to kill process: ", err)
	}
	// fmt.Printf("Waiting (%d)\n", runtime.NumGoroutine())

	// if err := cmd.Wait(); err != nil {
	//	log.Fatal("failed to wait on process: ", err)
	// }
	fmt.Printf("Done killing (%d)\n", runtime.NumGoroutine())
}

func chromeDPscreenshot(outPath string) {
	// create contex

	// ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithDebugf(log.Printf))
	// Create chrome instance
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Create a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// capture screenshot of an element
	var buf []byte
	if err := chromedp.Run(ctx, traceScreenshot(`http://127.0.0.1:32000/trace`, &buf)); err != nil {
		// if err := chromedp.Run(ctx, tracetest(`http://www.google.com`, &buf)); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile(outPath, buf, 0644); err != nil {
		log.Fatal(err)
	}

}
func tracetest(urlstr string, res *[]byte) chromedp.Tasks {
	// chromedp.Sleep(15 * time.Second),
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Screenshot("body", res, chromedp.NodeVisible, chromedp.ByQuery),
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
		// chromedp.SetAttributeValue(mouseModeSelector, "style.display", "none"),
		chromedp.Screenshot(traceViewerId, res, chromedp.NodeVisible, chromedp.ByID),
	}
}
