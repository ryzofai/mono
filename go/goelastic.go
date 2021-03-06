package main

import (
	"flag"
	"fmt"
	"log"
	"time"
	"context"

	"github.com/olivere/elastic"
)


type Logstream struct {
	Created  time.Time             `json:"created,omitempty"`
	Message  string                `json:"message"`
}


func main() {
	var (
		url   = flag.String("url", "http://192.168.201.215:9200", "Elasticsearch URL")
		sniff = flag.Bool("sniff", true, "Enable or disable sniffing")
		//index = 
	)
	flag.Parse()
	log.SetFlags(0)

	if *url == "" {
		*url = "http://127.0.0.1:9200"
	}

	// Create an Elasticsearch client
	client, err := elastic.NewClient(elastic.SetURL(*url), elastic.SetSniff(*sniff))
	if err != nil {
		log.Fatal(err)
	}
	_ = client

	// Just a status message
	fmt.Println("Connection succeeded")

	t := time.Now()
	fmt.Println()
	fmt.Println("logstreamer-1.0.0-" + string(t.Year()) + string(t.Month()) + string(t.Day()))

	/*_, err = client.CreateIndex("logstreamer-").Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}*/

	// Index a tweet (using JSON serialization)

	logster := Logstream{Created: time.Now(), Message: "test log 4"}
	put1, err := client.Index().
		Index("logstreamer-").
		Type("doc").
		//Id("1").
		BodyJson(logster).
		Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Indexed logs %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)

	// Flush to make sure the documents got written.
	_, err = client.Flush().Index("logstreamer").Do(context.Background())
	if err != nil {
		panic(err)
	}

	// Use the IndexExists service to check if a specified index exists.
	/*exists, err := client.IndexExists("logstream-").Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("logstream").Do(context.Background())
		if err != nil {
			// Handle error
			fmt.Println(err)
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
			fmt.Println("not ack")
		}
	}*/

}
