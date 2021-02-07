package main

import (
	"io/ioutil"
	"log"
	"os"
	"syscall"
)

type queryResult struct {
	name string
	info string
	url  string
}

// getLatestEpisode get the latest available episode of given name from given sources
// and push into channel.
// If no there's no available episode from given sources, it will push empty info
// into channel.
func getLatestEpisode(name string, sources []ComicSource, ch chan<- queryResult) {
	for _, source := range sources {
		if source.IsSupported(name) {
			data := source.GetLatestEpisode(name)
			ch <- queryResult{
				name: name,
				info: data,
				url:  source.GetURL(name),
			}
			return
		}
	}
	log.Printf("No available episode for %s\n", name)
	ch <- queryResult{
		name: name,
		info: "",
		url:  "",
	}
}

func getLatestComicEpisode(comicName string, ch chan<- queryResult) {
	sources := []ComicSource{&ComicBus{}}
	log.Printf("processing comicName = %+v\n", comicName)
	getLatestEpisode(comicName, sources, ch)
}

func getLatestAnimeEpisode(animeName string, ch chan<- queryResult) {
	sources := []ComicSource{&Iqiyi{}}
	log.Printf("processing animeName = %+v\n", animeName)
	getLatestEpisode(animeName, sources, ch)
}

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
		"Seven-deadly-sins",
		"Attack-on-Titan",
		"Demon-Slayer",
	}
	var animes = []string{
		"One-piece",
		"One-punch",
		"Demon-Slayer",
		"Two-Hit-Multi-Target-Attacks",
	}
	updateComics(jsonFile, comics, animes)
}
