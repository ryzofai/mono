package main

import (
"fmt"
"time"
"bufio"
"os"
"strings"
"gopkg.in/gomail.v2"
"github.com/tkanos/gonfig"
"flag"
"log"
"github.com/olivere/elastic"
"context"
"strconv"
"io"
)

type Configuration struct {
	Email_subject 		string
	Email_from 			string
	Email_to 			string
	Smtp_ip			 	string
	Smtp_port			int
	Files_to_monitor	string
	Dir_to_parse		string
	Alert_terms			string
	Elasticsearch_URL	string
	Elasticsearch_Index	string
}

type Logstream struct {
	Created  time.Time             `json:"created,omitempty"`
	Message  string                `json:"message"`
}

func main() {
	configuration := Configuration{}
	err := gonfig.GetConf("config.json", &configuration)
	if err != nil {  
		fmt.Println(err)
	}
	c := make(chan string)
	files := strings.Split(configuration.Files_to_monitor, ",")
	for i := range files {
		// race condition on logging.. =( 
		go prospector(files[i], c)

		sendTo := strings.Split(configuration.Email_to, ",")
		for email := range sendTo {
			go sendMail("MBA", "log streamer started", sendTo[email])
			// go sendMail("MBA", "log streamer started", configuration.Email_to)
		}
		time.Sleep(5 * time.Second)
	}
	x := <-c
	fmt.Println(x)
}

var counter int = 1

func prospector(files string, c chan string) {
	for {
		configuration := Configuration{}
		err := gonfig.GetConf("config.json", &configuration)
		if err != nil {  
			fmt.Println(err)
		}
		file, err := os.Open(files)
		if err != nil {
			fmt.Println(err)
		}
		line := bufio.NewScanner(file)
		for line.Scan() {
			alertTerms := strings.Split(configuration.Alert_terms, ",")
			for terms := range alertTerms {
				if strings.Contains(line.Text(), alertTerms[terms]) {
					sendTo := strings.Split(configuration.Email_to, ",")
					for email := range sendTo {
						go sendMail(files, line.Text(), sendTo[email])
						if checkMsg(files + "|" + strconv.Itoa(counter) + "|" + line.Text()) == false {
							elasticlogger(files + "|" + strconv.Itoa(counter) + "|" + line.Text()) // date(no timestamp), filename, save line number, specific log msg	
						}
					}
				}
			}
			counter++
			fmt.Println(files + " ||| " + strconv.Itoa(counter))
		}
		if err := line.Err(); err != nil {
			fmt.Println(err)
		}
		counter = 1
		time.Sleep(10 * time.Second)
	}
	c <- "running..."
}

func sendMail(filename string, message string, to string) {
	configuration := Configuration{}
	err := gonfig.GetConf("config.json", &configuration)
	if err != nil {  
		fmt.Println(err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", configuration.Email_from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", configuration.Email_subject)
	m.SetBody("text/text", "Alert term detected with the following details:\n" + "FILENAME: "+ filename + "" + "\nLINE CONTAINS: " + message)
	// m.SetHeader("To", result[i])
	// m.SetAddressHeader("Cc", "", "")
	// m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer(configuration.Smtp_ip, configuration.Smtp_port, "", "")

	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
	}
	time.Sleep(2 * time.Second)
}

var (
	url   = flag.String("url", "http://192.168.201.215:9200", "Elasticsearch URL")
	sniff = flag.Bool("sniff", true, "Enable or disable sniffing")
	)

func elasticlogger(logthis string) {
	configuration := Configuration{}
	err := gonfig.GetConf("config.json", &configuration)
	if err != nil {  
		fmt.Println(err)
	}

	flag.Parse()
	log.SetFlags(0)

	if *url == "" {
		fmt.Println("url null")
		return
	}
	// Create an Elasticsearch client - *todo: update this to not use pointers to allow configuration from config file
	client, err := elastic.NewClient(elastic.SetURL(*url), elastic.SetSniff(*sniff))
	if err != nil {
		fmt.Println(err)
	}
	_ = client
	// print status message
	fmt.Println("Connection succeeded")

	//fmt.Println(configuration.Elasticsearch_Index + string(t.Year()) + string(t.Month()) + string(t.Day()))
	var msg string = "\"" + logthis + "\""
	//t := time.Now()
	//var indexDate = string(t.Year()) + string(t.Month()) + string(t.Day())
	logit := Logstream{Created: time.Now(), Message: msg}
	put1, err := client.Index().
	Index(configuration.Elasticsearch_Index).
	Type("doc").
		//Id("1").
	BodyJson(logit).
	Do(context.Background())
	if err != nil {
		// Handle error
		fmt.Println(err)
	}
	fmt.Printf("Indexed logs %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)

	// Flush to make sure the documents got written.
	_, err = client.Flush().Index(configuration.Elasticsearch_Index).Do(context.Background())
	if err != nil {
		fmt.Println(err)
	} else {
		logMsg(logthis)
	}
}

func logMsg(logstr string) {
	f, err := os.OpenFile("log.txt", os.O_APPEND, 0600)
	defer f.Close()
	if fileExists("log.txt") != false {
		if _, err = f.WriteString(logstr + "\n")
		err != nil {
			panic(err)
		}
		f.Close()
	} else {
		nfile, err := os.Create("log.txt")
		if err != nil {
			log.Fatal("Cannot create file", err)
		}
		nfile.Close()
	}
	f.Close()
}

func checkMsg(msg string) bool {
	if fileExists("log.txt") == true {
		file, err := os.Open("log.txt")
		if err != nil {
			fmt.Println(err)
		}
		line := bufio.NewScanner(file)
		for line.Scan() {
			if line.Text() == msg {
				return true
			}
		}	
		file.Close()	
	} else {
		efile, err := os.Create("log.txt")
		if err != nil {
			log.Fatal("Cannot create file", err)
		}
		//fmt.Println("in elsefier")
		//fmt.Fprintf(efile, msg + "\n")
		efile.Close()
		return false
		// create file then log msg
		// if no errors return true
		// else return false
	}
	return false
}

func fileExists(fname string) bool {
	if _, err := os.Stat(fname); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
// add go routine function log maintenance.. clear log file contents daily
// rename file only


func renameFile() {
	err :=  os.Rename("log.txt", "log.x.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func deleteFile() {
	// delete file
	var err = os.Remove("log.txt")
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}

func replicator() {
	from, err := os.Open("log.w.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer from.Close()

	to, err := os.OpenFile("log.r.txt", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		log.Fatal(err)
	}
}
