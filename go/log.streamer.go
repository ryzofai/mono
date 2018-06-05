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
"sync"
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


var wg sync.WaitGroup
func main() {
	//progArgs := os.Args
	//fmt.Println(progArgs)
	//fmt.Println(progArgs[1])
	configuration := Configuration{}
	err := gonfig.GetConf("config.cfg", &configuration)
	if err != nil {  
		fmt.Println(err)
	}
    var wg sync.WaitGroup
	for {
		files := strings.Split(configuration.Files_to_monitor, ",")
		for i := range files {
			wg.Add(1)
			go prospector(&wg, files[i])
			time.Sleep(2 * time.Second)
		}
		wg.Wait()
		fmt.Println("end...")
		time.Sleep(30 * time.Second)
	}
}

func prospector(wg *sync.WaitGroup, files string) {
	var counter int = 1
	defer wg.Done()
	configuration := Configuration{}
	err := gonfig.GetConf("config.cfg", &configuration)
	if err != nil {  
		fmt.Println(err)
	}
	file, err := os.Open(files)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
	}
	line := bufio.NewScanner(file)
	fmt.Println(file)
	for line.Scan() {
		alertTerms := strings.Split(configuration.Alert_terms, ",")
		for terms := range alertTerms {
			if strings.Contains(line.Text(), alertTerms[terms]) {
				fmt.Println(alertTerms[terms])
				fmt.Println(line.Text())
				sendTo := strings.Split(configuration.Email_to, ",")
				for email := range sendTo {
					sendMail(files, line.Text(), sendTo[email])
					if checkMsg(files + "|" + strconv.Itoa(counter) + "|" + line.Text()) == false {
						fmt.Println(line.Text())
						fmt.Println(counter)
						// add go channel to not terminate prematurely 
						elasticlogger(files + "|" + strconv.Itoa(counter) + "|" + line.Text()) // date(no timestamp), filename, save line number, specific log msg	
					}
				}
			}
		}
		counter++
	}
}

func sendMail(filename string, message string, to string) {
	configuration := Configuration{}
	err := gonfig.GetConf("config.cfg", &configuration)
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
	time.Sleep(1 * time.Second)
}

var (
	url   = flag.String("url", "http://192.168.201.215:9200", "Elasticsearch URL")
	sniff = flag.Bool("sniff", true, "Enable or disable sniffing")
	)

func elasticlogger(logthis string) {
	configuration := Configuration{}
	err := gonfig.GetConf("config.cfg", &configuration)
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
	//c <- "running..."
}

func logMsg(logstr string) {
	f, err := os.OpenFile("log.txt", os.O_APPEND, 0600)
//defer f.Close()
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
				fmt.Println("found in file: " + msg)
				return true
			}
		}	
		file.Close()	
	} else {
		efile, err := os.Create("log.txt")
		if err != nil {
			log.Fatal("Cannot create file", err)
		}
		//fmt.Println("in elsefier: " + msg)
		//fmt.Fprintf(efile, msg + "\n")

		efile.Close()
		checkMsg(msg)
		//return false
	// create file then log msg
	// if no errors return true
	// else return false
	}
	//fmt.Println("found in file: " + msg)
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
	//defer from.Close()

	to, err := os.OpenFile("log.r.txt", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	//defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		log.Fatal(err)
	}
}
