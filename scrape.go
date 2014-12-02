package main

import (
	"log"

	"github.com/jaisingh/mls-scraper/mls"
)

func main() {
	log.SetFlags(log.Lshortfile + log.Ltime)
	log.Println("Starting...")

	//var err error
	if err := mls.Init(); err != nil {
		log.Fatal(err)
	}

	/*
		m, _ := mls.GetMLSIDs()
			for _, l := range m {
				i, _ := mls.GetListing(l)
				log.Printf("%#v\n", i)
			}*/

	l, _ := mls.GetListing("71734795")
	log.Printf("%#v\n", l)
}
