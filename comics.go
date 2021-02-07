package main

import (
	"encoding/json"
	"fmt"
)

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
			fmt.Printf("Got new update: episode = %+v (%s)\n", queryRes.info, queryRes.url)
			shouldUpdate = true
		}
	}
	for i := 0; i < len(animes); i++ {
		queryRes := <-animeChan
		oriEps, ok := comic.Anime[queryRes.name]
		if !ok || oriEps != queryRes.info {
			comic.Anime[queryRes.name] = queryRes.info
			fmt.Printf("Got new update: episode = %+v (%s)\n", queryRes.info, queryRes.url)
			shouldUpdate = true
		}
	}
	return
}
