package main

import "testing"

func TestComicbusIsSupportedName(t *testing.T) {
	src := ComicBus{}
	comicBusList = map[string]string{}
	src.AddList("one-piece", "103")
	// Comic name check is case insensitive
	if !src.IsSupported("one-piece") {
		t.Errorf("expecting One-Piece to be supported, get false")
	}
	if !src.IsSupported("One-Piece") {
		t.Errorf("expecting One-Piece to be supported, get false")
	}
}

func TestComicbusNotIsSupportedName(t *testing.T) {
	src := ComicBus{}
	iqiyiList = map[string]string{}
	if src.IsSupported("temp") {
		t.Errorf("Expect temp no to be supported")
	}
	src.AddList("temp", "test")
	if !src.IsSupported("temp") {
		t.Errorf("Expect temp to be supported")
	}
}

func TestComicbusGetURL(t *testing.T) {
	src := ComicBus{}
	src.AddList("temp", "test")
	res := src.GetURL("temp")
	exp := "https://www.comicbus.com/html/test.html"
	if res != exp {
		t.Errorf("Expect %s, got %s", exp, res)
	}
}
