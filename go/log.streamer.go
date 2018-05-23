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
		fmt.Println(out)
		if string(out) == "" {
			fmt.Println("orodsdwfwoh")
		}
		/*if strings.Contains(string(out), "process not found") {
			fmt.Println("orodsdwfwoh")
		}*/

		if err != nil {
			fmt.Println("error not nill")
			if ctr >= 5 {
				exec.Command("cmd", "/C", "taskkill /IM loggit.exe /F")
				fmt.Println(string(out))
				fmt.Println("taskill")
				ctr = 0
				fmt.Println("counter:")
				fmt.Println(ctr)
			}

			//t := strings.Split(out, "\n")
			/*for i := range out {
				ctr2++
				fmt.Println(i)
				fmt.Println(strconv.Itoa(ctr2))
				}*/
				ctr++
			} 
			//fmt.Println("error nill")
			cmnd2 := exec.Command("loggit.exe", "arg")
			cmnd2.Start()
			//fmt.Println("process not found " + strconv.Itoa(ctr))
			//ctr = 0



			
			time.Sleep(10 * time.Second)
		}
	}
