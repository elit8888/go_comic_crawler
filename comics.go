package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type queryResult struct {
	name string
	info string
	url  string
}

// Comics struct contains latest parsed episodes
type Comics struct {
	Comic map[string]string `json:"comic"`
	Anime map[string]string `json:"anime"`
}

// FromJSON restore the data from given json byte array
func (comic *Comics) FromJSON(data []byte) error {
	err := json.Unmarshal(data, comic)
	if err != nil {
		return err
	}
	if comic.Comic == nil {
		comic.Comic = make(map[string]string)
	}
	if comic.Anime == nil {
		comic.Anime = make(map[string]string)
	}
	return nil
}

// ToJSON dump data to json byte array
func (comic *Comics) ToJSON() ([]byte, error) {
	return json.MarshalIndent(comic, "", "  ")
}

// UpdateEpisodes updates the given comics and animes,
// shouldUpdate indicates whether there's any newer episode available.
func (comic *Comics) UpdateEpisodes(comics []string, animes []string) (shouldUpdate bool, err error) {
	comicChan, animeChan := make(chan queryResult), make(chan queryResult)
	for _, comicName := range comics {
		go getLatestComicEpisode(comicName, comicChan)
	}
	for _, animeName := range animes {
		go getLatestAnimeEpisode(animeName, animeChan)
	}

	// receive update
	for i := 0; i < len(comics); i++ {
		queryRes := <-comicChan
		oriEps, ok := comic.Comic[queryRes.name]
		if !ok || oriEps != queryRes.info {
			comic.Comic[queryRes.name] = queryRes.info
			fmt.Printf("%s got new episode: %+v (%s)\n", queryRes.name, queryRes.info, queryRes.url)
			shouldUpdate = true
		}
	}
	for i := 0; i < len(animes); i++ {
		queryRes := <-animeChan
		oriEps, ok := comic.Anime[queryRes.name]
		if !ok || oriEps != queryRes.info {
			comic.Anime[queryRes.name] = queryRes.info
			fmt.Printf("%s got new episode: %+v (%s)\n", queryRes.name, queryRes.info, queryRes.url)
			shouldUpdate = true
		}
	}
	return
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
