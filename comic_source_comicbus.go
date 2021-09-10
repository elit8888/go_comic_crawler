package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
)

// ComicBus implements ComicSource, which query episodes from comicbus website.
type ComicBus struct{}

var comicBusList = map[string]string{
	"one-piece":       "103",
	"one-punch":       "10406",
	"attack-on-titan": "7340",
	"jujutsu-kaisen":  "15790",
}

var comicbusBaseURL = "https://www.comicbus.com/html/%s.html"

// IsSupported returns whether the given comic name is supported.
func (comicbus *ComicBus) IsSupported(name string) bool {
	if _, ok := comicBusList[strings.ToLower(name)]; ok {
		return true
	}
	return false
}

// GetLatestEpisode query latest episode from iqiyi website
func (comicbus *ComicBus) GetLatestEpisode(name string) string {
	var data string
	url := comicbus.GetURL(name)
	log.Printf("Query latest episode of %s from %s", name, url)
	collector := colly.NewCollector()
	collector.OnHTML("#Comic", func(element *colly.HTMLElement) {
		data = strings.Split(strings.Fields(element.Text)[2], "-")[1]
		log.Printf("Comic %v got %v\n", name, data)
	})
	collector.Visit(url)
	return data
}

// GetURL returns the string of the url to be queried.
func (comicbus *ComicBus) GetURL(name string) string {
	name = strings.ToLower(name)
	return fmt.Sprintf(comicbusBaseURL, comicBusList[name])
}

// AddList add a mapping entry for given name, and it overwrite the original one
// if entry exists.
func (comicbus *ComicBus) AddList(name, page string) {
	comicBusList[name] = page
}
