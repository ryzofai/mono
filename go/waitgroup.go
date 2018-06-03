package main
import (
"fmt"
"os"
"strings"
"bufio"
"sync"
)

func main() {
	var wg sync.WaitGroup
	for {
		files := strings.Split("C:\\Projects\\Go\\file1.txt,C:\\Projects\\Go\\file2.txt", ",")
		for file := range files {
			wg.Add(1)
			go readfile(&wg, files[file])
		}
		wg.Wait()
		fmt.Println("end...")
		break
	}
}

func readfile(wg *sync.WaitGroup,path string) {
	defer wg.Done()
	inFile, _ := os.Open(path)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines) 
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
