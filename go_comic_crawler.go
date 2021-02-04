package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/gocolly/colly/v2"
)

// Comics struct contains latest parsed episodes
type Comics struct {
	Comic map[string]string `json:"comic"`
	Anime map[string]string `json:"anime"`
}

type queryResult struct {
	name string
	info string
}

type urls map[string]string

func getLatestComicURL(comicName string, comicURL string, res chan<- queryResult) {
	log.Printf("processing comicName = %+v\n", comicName)
	collector := colly.NewCollector()
	data := ""
	collector.OnHTML("#Comic", func(element *colly.HTMLElement) {
		data = strings.Split(strings.Split(element.Text, " ")[1], "-")[1]
		log.Printf("Comic %v got %v\n", comicName, data)
	})
	collector.Visit(comicURL)
	res <- queryResult{
		name: comicName,
		info: data,
	}
}

func getLatestAnimeURL(animeName string, animeURL string, res chan<- queryResult) {
	log.Printf("processing animeName = %+v\n", animeName)
	collector := colly.NewCollector()
	data := ""
	collector.OnHTML("span[class=slide-tag]", func(element *colly.HTMLElement) {
		data = element.Text
		log.Printf("Anime %v got %v\n", animeName, data)
	})
	collector.Visit(animeURL)
	res <- queryResult{
		name: animeName,
		info: data,
	}
}

// FromJSON restore the data from given json byte array
func (comics *Comics) FromJSON(data []byte) error {
	err := json.Unmarshal(data, comics)
	if err != nil {
		return err
	}
	if comics.Comic == nil {
		comics.Comic = make(map[string]string)
	}
	if comics.Anime == nil {
		comics.Anime = make(map[string]string)
	}
	return nil
}

// ToJSON dump data to json byte array
func (comics *Comics) ToJSON() ([]byte, error) {
	return json.MarshalIndent(comics, "", "  ")
}

// UpdateEpisodes updates the given comics and animes,
// shouldUpdate indicates whether there's any newer episode available.
func (comics *Comics) UpdateEpisodes(comicURLs urls, animeURLs urls) (shouldUpdate bool, err error) {
	comicChan, animeChan := make(chan queryResult), make(chan queryResult)
	for comicName, comicURL := range comicURLs {
		go getLatestComicURL(comicName, comicURL, comicChan)
	}
	for animeName, animeURL := range animeURLs {
		go getLatestAnimeURL(animeName, animeURL, animeChan)
	}

	// receive update
	for i := 0; i < len(comicURLs); i++ {
		queryRes := <-comicChan
		val, ok := comics.Comic[queryRes.name]
		if !ok || val != queryRes.info {
			comics.Comic[queryRes.name] = queryRes.info
			log.Printf("Got new update: episode = %+v\n", queryRes.info)
			shouldUpdate = true
		}
	}
	for i := 0; i < len(animeURLs); i++ {
		queryRes := <-animeChan
		val, ok := comics.Anime[queryRes.name]
		if !ok || val != queryRes.info {
			comics.Anime[queryRes.name] = queryRes.info
			log.Printf("Got new update: episode = %+v\n", queryRes.info)
			shouldUpdate = true
		}
	}
	return
}

func updateComics(jsonFile string, comicURLs urls, animeURLs urls) {
	var comics Comics
	byteValue, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		if e, ok := err.(*os.PathError); ok && e.Err == syscall.ENOENT {
			log.Println("File not exists, instantiate a new data")
			byteValue = []byte("{}") // empty json data
		} else {
			log.Fatalf("Cannot read json from %s, error = %+v\n", jsonFile, e)
		}
	}
	if err = comics.FromJSON(byteValue); err != nil {
		log.Fatalln("Initiate data failed, error:", err)
	}

	// Parse website to update data
	shouldUpdate, err := comics.UpdateEpisodes(comicURLs, animeURLs)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("record = %+v\n", comics)
	content, err := comics.ToJSON()
	if err != nil {
		log.Fatal(err)
	}

	// save record if there's any update
	if shouldUpdate {
		log.Printf("Write json content to %s\n", jsonFile)
		ioutil.WriteFile(jsonFile, content, 0644)
	} else {
		log.Println("No update available")
	}
}

func main() {
	jsonFile := "crawl_record.json"
	var comicURLs = urls{
		"One-piece":         "https://www.comicbus.com/html/103.html",
		"One-punch":         "https://www.comicbus.com/html/10406.html",
		"Seven-deadly-sins": "https://www.comicbus.com/html/9418.html",
		"Attack-on-Titan":   "https://www.comicbus.com/html/7340.html",
		"Demon-Slayer":      "https://www.comicbus.com/html/14132.html",
	}
	var animeURLs = urls{
		"One-piece":                    "https://tw.iqiyi.com/a_19rrh8ngb1.html",
		"One-punch":                    "https://tw.iqiyi.com/a_19rrhtbgxd.html",
		"Demon-Slayer":                 "https://tw.iqiyi.com/a_19rrhrnr05.html",
		"Two-Hit-Multi-Target-Attacks": "https://tw.iqiyi.com/a_19rri0oj3l.html",
	}
	updateComics(jsonFile, comicURLs, animeURLs)
}
