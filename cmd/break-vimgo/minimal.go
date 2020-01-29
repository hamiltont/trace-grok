package main

import (
	"os/exec"
	"time"
)

func main() {

	cmd := exec.Command("/bin/sh", "-c", "watch date > date.txt")
	// cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	time.AfterFunc(3*time.Second, func() { cmd.Process.Kill() })
	// time.AfterFunc(3*time.Second, func() {
	// 	syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	// })
	cmd.Run()

}
