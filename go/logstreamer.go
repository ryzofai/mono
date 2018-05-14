package main

import(
"fmt"
"os"
"io/ioutil"
"log"
"time"
"runtime"
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
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("here1")
	cx := make(chan string)
	for _, f := range files {
		fileInfo, err := os.Stat(f.Name())
		if err != nil {
			fmt.Println(err)
		}
		t := fileInfo.ModTime()
		n := time.Now()
		
		if t.Day() == n.Day() {
			// check filename if in prospector list before running:
			go prospector(f.Name(), cx)
		}
		fmt.Println("here2")
	}
	fmt.Println(runtime.NumGoroutine())
	time.Sleep(5 * time.Second)
	c <- cx
}

// attach agent to file for continuous monitoring
func prospector(s string, c chan string) {
	for {
		time.Sleep(5 * time.Second)
		fmt.Println("prospector: " + s)
		fmt.Println("here3")
		//cf_slice = append(cf_slice, s)
		//fmt.Println(cf_slice)
		//return // return if mod date !=curr date
	}
	c <- s
}

/*func launcher(c chan string) {
	if err != nil {
		log.Fatal(err)
	}
}*/


/*import(
	"fmt"
	"reflect"
)

func main() {
	items := []int{1,2,3,4,5,6}
	fmt.Println(SliceExists(items, 5)) // returns true
	fmt.Println(SliceExists(items, 10)) // returns false
}

func SliceExists(slice interface{}, item interface{}) bool {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		panic("SliceExists() given a non-slice type")
	}

	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Interface() == item {
			return true
		}
	}

	return false
}*/
