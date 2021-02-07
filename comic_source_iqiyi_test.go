package main

import (
	"testing"
)

func TestIqiyiIsSupportedName(t *testing.T) {
	src := Iqiyi{}
	iqiyiList = map[string]string{}
	src.AddList("one-piece", "a_19rrh8ngb1")
	// Comic name check is case insensitive
	if !src.IsSupported("one-piece") {
		t.Errorf("expecting One-Piece to be supported, get false")
	}
	if !src.IsSupported("One-Piece") {
		t.Errorf("expecting One-Piece to be supported, get false")
	}
}

func TestIqiyiNotIsSupportedName(t *testing.T) {
	src := Iqiyi{}
	iqiyiList = map[string]string{}
	if src.IsSupported("temp") {
		t.Errorf("Expect temp no to be supported")
	}
	src.AddList("temp", "test")
	if !src.IsSupported("temp") {
		t.Errorf("Expect temp to be supported")
	}
}

func TestIqiyiGetURL(t *testing.T) {
	src := Iqiyi{}
	src.AddList("temp", "test")
	res := src.GetURL("temp")
	exp := "https://tw.iqiyi.com/test.html"
	if res != exp {
		t.Errorf("Expect %s, got %s", exp, res)
	}
}
