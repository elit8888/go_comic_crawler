package main

// ComicSource is the interface that is used to query comic result, or get url.
//
// IsSupported check if the given comic name is supported by the source.
// GetLatestEpisode query the latest available episode from the source.
// GetURL returns the endpoint url for the given comic name.
type ComicSource interface {
	IsSupported(name string) bool
	GetLatestEpisode(name string) string
	GetURL(name string) string
}
