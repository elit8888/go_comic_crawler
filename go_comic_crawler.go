package main

import (
	"io/ioutil"
	"log"
	"os"
	"syscall"
)

func updateComics(jsonFile string, comics []string, animes []string) {
	var comic Comics

	// Read data from file
	byteValue, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		if e, ok := err.(*os.PathError); ok && e.Err == syscall.ENOENT {
			log.Println("File not exists, instantiate a new data")
			byteValue = []byte("{}") // empty json data
		} else {
			log.Fatalf("Cannot read json from %s, error = %+v\n", jsonFile, e)
		}
	}
	if err = comic.FromJSON(byteValue); err != nil {
		log.Fatalln("Initiate data failed, error:", err)
	}

	// Parse website to update data
	shouldUpdate, err := comic.UpdateEpisodes(comics, animes)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("record = %+v\n", comic)

	// Save record if there's any update
	content, err := comic.ToJSON()
	if err != nil {
		log.Fatal(err)
	}

	if shouldUpdate {
		log.Printf("Write json content to %s\n", jsonFile)
		ioutil.WriteFile(jsonFile, content, 0644)
	} else {
		log.Println("No update available")
	}
}

func main() {
	jsonFile := "crawl_record.json"
	var comics = []string{
		"One-piece",
		"One-punch",
		"Attack-on-Titan",
		"Jujutsu-Kaisen",
	}
	var animes = []string{
		"One-piece",
		"One-punch",
		"Two-Hit-Multi-Target-Attacks",
	}
	updateComics(jsonFile, comics, animes)
}
