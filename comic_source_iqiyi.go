package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
)

// Iqiyi implements ComicSource, which query episodes from iqiyi website.
type Iqiyi struct{}

var iqiyiList = map[string]string{
	"one-piece":                    "a_19rrh8ngb1",
	"one-punch":                    "a_19rrhtbgxd",
	"two-hit-multi-target-attacks": "a_19rri0oj3l",
}

var iqiyiBaseURL = "https://tw.iqiyi.com/%s.html"

// IsSupported returns whether the given comic name is supported.
func (iqiyi *Iqiyi) IsSupported(name string) bool {
	if _, ok := iqiyiList[strings.ToLower(name)]; ok {
		return true
	}
	return false
}

// GetLatestEpisode query latest episode from iqiyi website
func (iqiyi *Iqiyi) GetLatestEpisode(name string) string {
	var data string
	url := iqiyi.GetURL(name)
	log.Printf("Query latest episode of %s from %s", name, url)
	collector := colly.NewCollector()
	collector.OnHTML("span[class=slide-tag]", func(element *colly.HTMLElement) {
		data = element.Text
		log.Printf("Anime %v got %v\n", name, data)
	})
	collector.Visit(url)
	return data
}

// GetURL returns the string of the url to be queried.
func (iqiyi *Iqiyi) GetURL(name string) string {
	name = strings.ToLower(name)
	return fmt.Sprintf(iqiyiBaseURL, iqiyiList[name])
}

// AddList add a mapping entry for given name, and it overwrite the original one
// if entry exists.
func (iqiyi *Iqiyi) AddList(name, page string) {
	iqiyiList[name] = page
}
