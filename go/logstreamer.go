package main

import(
"fmt"
"os"
"io/ioutil"
"log"
"time"
)

//var cf_slice []string // array of files in directory
//var cf_clear []string // clear files in list

func main() {
	c := make(chan string)
	go inspector(c)
	x := <-c
	fmt.Println(x)
}

// inspect dir for changes in mod date every 10 secs
func inspector(c chan string) {
	// if file mod date == curr date launch prospector	
	for {
		files, err := ioutil.ReadDir("./")
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range files {

			fileInfo, err := os.Stat(f.Name())
			if err != nil {
				fmt.Println(err)
			}
			t := fileInfo.ModTime()
			n := time.Now()
			c := make(chan string)
			if t.Day() == n.Day() {
				//cf = append(cf, f.Name())
				//fmt.Println(cf)
				go prospector(f.Name(), c)
				// x := <-c
				// fmt.Println(x)
			}
		}
		time.Sleep(10 * time.Second)
	}
	c <- "here"
}

// attach agent to file for continuous monitoring
func prospector(s string, c chan string) {
	for {
		time.Sleep(10 * time.Second)
		fmt.Println("prospector: " + s)
		//cf_slice = append(cf_slice, s)
		//fmt.Println(cf_slice)
		return // return if mod date !=curr date
	}
	c <- s
}

/*func launcher(c chan string) {


	if err != nil {
		log.Fatal(err)
	}
}*/
