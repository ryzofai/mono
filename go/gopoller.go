package main

import (
    "fmt"
	"time"
	"bufio"
	"log"
	"os"
	// "net/smtp"
	"strings"
	"gopkg.in/gomail.v2"
	"github.com/tkanos/gonfig"
	//"github.com/shomali11/util/xconcurrency"
	// "strconv"
	// "context"
	// "github.com/segmentio/kafka-go"
)

type Configuration struct {
	Email_subject 		string
	Email_from 			string
	Email_to 			string
	Smtp_ip			 	string
	Smtp_port			int
	Files_to_parse		string
}

func main() {
	configuration := Configuration{}
	err := gonfig.GetConf("config.json", &configuration)
	if err != nil {  
		log.Fatal(err)
	}
	c := make(chan string)
	files := strings.Split(configuration.Files_to_parse, ",")
	//fmt.Println(files[0])
	//fmt.Println(files[1])
	//xconcurrency.Parallelize(prospector("note.txt"), prospector("note2.txt"))
	for i := range files {
		var f = string(files[i])
		//fmt.Println(f)
		//xconcurrency.Parallelize(prospector(f), func2)
		go prospector(f, c)
		//time.Sleep(10000 * time.Millisecond)
	}
	x, y := <-c, <-c
	fmt.Println(x, y, x+y)
}

func prospector(files string, c chan string) {
	//fmt.Println("file: " + files)
	//c <- files + " test"
	configuration := Configuration{}
	err := gonfig.GetConf("config.json", &configuration)
	if err != nil {  
		log.Fatal(err)
	}
	// fmt.Println(files)
    file, err := os.Open(files)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// fmt.Println(configuration.files_to_parse);
	for {
		scanner := bufio.NewScanner(file)

		// fmt.Println("file in scanner: " + scanner.Text())
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "error") {
				// for each in send list
				result := strings.Split(configuration.Email_to, ",")
				for i := range result {
					// fmt.Println(result[i])
					goSendMail(files, scanner.Text(), result[i])
					
				}
			}
			// fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		time.Sleep(10000 * time.Millisecond)
	}
	c <- "running..."
}

func goSendMail(filename string, message string, to string) {

	configuration := Configuration{}
	err := gonfig.GetConf("config.json", &configuration)
	if err != nil {  
		log.Fatal(err)
	}
	/*fmt.Println(configuration.Email_from);
	fmt.Println(configuration.Email_to);
	fmt.Println(configuration.Email_subject);
	fmt.Println(configuration.Smtp_ip);
	fmt.Println(configuration.Smtp_port);*/

	m := gomail.NewMessage()
	m.SetHeader("From", configuration.Email_from)
	/*result := strings.Split(configuration.Email_to, ",")
	for i := range result {
		fmt.Println(result[i])
		m.SetHeader("To", result[i])
	}*/
	m.SetHeader("To", to)
	// m.SetAddressHeader("Cc", "", "")
	m.SetHeader("Subject", "Test!")
	m.SetBody("text/text", filename + ": " + message)
	// m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer(configuration.Smtp_ip, configuration.Smtp_port, "", "")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
