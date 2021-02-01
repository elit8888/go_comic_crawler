package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"

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

func readJSON(filename string) (Comics, error) {
	var comics Comics
	byteValue, err := ioutil.ReadFile(filename)
	if err != nil {
		return comics, err
	}
	json.Unmarshal(byteValue, &comics)
	return comics, nil
}

func writeJSON(filename string, comics Comics) error {
	jsonString, err := json.MarshalIndent(comics, "", "  ")
	if err != nil {
		log.Fatalf("Cannot write data to %v: %v", filename, err)
	}
	ioutil.WriteFile(filename, jsonString, 0644)
	return nil
}

func getLatestComicURL(comicName string, comicURL string, res chan<- queryResult) {
	log.Printf("processing comicName = %+v\n", comicName)
	collector := colly.NewCollector()
	data := ""
	collector.OnHTML("#Comic", func(element *colly.HTMLElement) {
		data = strings.Split(strings.Split(element.Text, " ")[1], "-")[1]
		log.Printf("Comic %v got %v", comicName, data)
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
		log.Printf("Anime %v got %v", animeName, data)
	})
	collector.Visit(animeURL)
	res <- queryResult{
		name: animeName,
		info: data,
	}
}

func updateComics(
	jsonFile string,
	comicURLs map[string]string,
	animeURLs map[string]string,
) {
	record, err := readJSON(jsonFile)
	if err != nil {
		log.Printf("Cannot data from %v: %s", jsonFile, err)
		log.Printf("Begin as empty data")
	}
	log.Printf("%+v\n", record)
	if record.Comic == nil {
		record.Comic = make(map[string]string)
	}
	if record.Anime == nil {
		record.Anime = make(map[string]string)
	}

	// Parse website to update data
	comicChan, animeChan := make(chan queryResult), make(chan queryResult)
	for comicName, comicURL := range comicURLs {
		go getLatestComicURL(comicName, comicURL, comicChan)
	}
	for animeName, animeURL := range animeURLs {
		go getLatestAnimeURL(animeName, animeURL, animeChan)
	}

	// receive update
	shouldUpdate := false
	for i := 0; i < len(comicURLs); i++ {
		queryRes := <-comicChan
		val, ok := record.Comic[queryRes.name]
		if !ok || val != queryRes.info {
			record.Comic[queryRes.name] = queryRes.info
			shouldUpdate = true
		}
	}
	for i := 0; i < len(animeURLs); i++ {
		queryRes := <-animeChan
		val, ok := record.Anime[queryRes.name]
		if !ok || val != queryRes.info {
			record.Anime[queryRes.name] = queryRes.info
			shouldUpdate = true
		}
	}
	log.Printf("record = %+v\n", record)

	// save record if there's any update
	if shouldUpdate {
		log.Printf("Update json file")
		writeJSON(jsonFile, record)
	} else {
		log.Printf("No update available")
	}
}

func main() {
	jsonFile := "crawl_record.json"
	var comicURLs = map[string]string{
		"One-piece":         "https://www.comicbus.com/html/103.html",
		"One-punch":         "https://www.comicbus.com/html/10406.html",
		"Seven-deadly-sins": "https://www.comicbus.com/html/9418.html",
		"Attack-on-Titan":   "https://www.comicbus.com/html/7340.html",
		"Demon-Slayer":      "https://www.comicbus.com/html/14132.html",
	}
	var animeURLs = map[string]string{
		"One-piece":                    "https://tw.iqiyi.com/a_19rrh8ngb1.html",
		"One-punch":                    "https://tw.iqiyi.com/a_19rrhtbgxd.html",
		"Demon-Slayer":                 "https://tw.iqiyi.com/a_19rrhrnr05.html",
		"Two-Hit-Multi-Target-Attacks": "https://tw.iqiyi.com/a_19rri0oj3l.html",
	}

	updateComics(jsonFile, comicURLs, animeURLs)
}
