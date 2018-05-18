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
		
		go prospector(files[i], c)
		sendTo := strings.Split(configuration.Email_to, ",")
		for email := range sendTo {
			go sendMail("MBA", "log streamer started", sendTo[email])
			// go sendMail("MBA", "log streamer started", configuration.Email_to)
		}
	}
	x := <-c
	fmt.Println(x)
}

func prospector(files string, c chan string) {
	configuration := Configuration{}
	err := gonfig.GetConf("config.json", &configuration)
	if err != nil {  
		fmt.Println(err)
	}
	file, err := os.Open(files)
	if err != nil {
		fmt.Println(err)
	}
	//defer file.Close()
	// rewrite: read only if modified date & time != saved modified date & time
	// 			read only from saved record line 
	var counter int

	for {
		line := bufio.NewScanner(file)
		for line.Scan() {
			alertTerms := strings.Split(configuration.Alert_terms, ",")
			for terms := range alertTerms {
				if strings.Contains(line.Text(), alertTerms[terms]) {
					sendTo := strings.Split(configuration.Email_to, ",")
					for email := range sendTo {
						go sendMail(files, line.Text(), sendTo[email])
						go elasticlogger(line.Text())
					}
				}
			}
		}
		if err := line.Err(); err != nil {
			fmt.Println(err)
		}
		//file.Close()
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
	// Create an Elasticsearch client
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
	}
	
}
