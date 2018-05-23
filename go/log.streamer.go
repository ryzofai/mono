package main

import (
//"log"
"os/exec"
"time"
"fmt"
//"strconv"
//"strings"
)

var ctr int

func main() {
	for {
		cmnd1 := exec.Command("cmd", "/C","tasklist|findstr loggit.exe")
		out, err := cmnd1.Output()
		cmnd1.Start()
		//fmt.Println(out)
		if string(out) == "" {
			cmnd2 := exec.Command("loggit.exe", "arg")
			cmnd2.Start()
		} else if ctr >= 5 {
			exec.Command("cmd", "/C", "taskkill /IM loggit.exe /F")
			fmt.Println("taskill " + string(out) + string(ctr))
		} else {
			fmt.Println(ctr)
			ctr++
		}
		if err != nil {}
		time.Sleep(10 * time.Second)
	}
}
