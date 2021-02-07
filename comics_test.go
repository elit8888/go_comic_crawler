package main

import "testing"

func TestComicFromJsonEmptyInput(t *testing.T) {
	var comic Comics
	err := comic.FromJSON([]byte(""))
	if err == nil {
		t.Errorf("Expect error from empty input")
	}
}

func TestComicFromJsonInvalidInput(t *testing.T) {
	var comic Comics
	err := comic.FromJSON([]byte("test"))
	if err == nil {
		t.Errorf("Expect error from unstructured input format")
	}
}

func TestEmptyInput(t *testing.T) {
	emptyInput := []string{
		"{\"comic\":{}}",
		"{\"anime\":{}}",
		"{\"anime\":{}, \"comic\":{}}",
		"{\"anime\":{}, \"comic\":{},\"other\":{}}",
		"{\"anime\":null, \"comic\":null}",
	}
	for _, test := range emptyInput {
		t.Run(test, func(t *testing.T) {
			var comic Comics
			err := comic.FromJSON([]byte(test))
			if err != nil {
				t.Errorf("Expect no error, got %+v", err)
			}
			if l := len(comic.Anime); l != 0 {
				t.Errorf("Expect no content in comic.Anime, got %d", l)
			}
			if l := len(comic.Comic); l != 0 {
				t.Errorf("Expect no content in comic.Comic, got %d", l)
			}
		})
	}
}

func TestFromJSON(t *testing.T) {
	var comic Comics
	err := comic.FromJSON([]byte("{\"comic\":{\"test1\":\"1\"}}"))
	if err != nil {
		t.Errorf("Expect no error, got error: %+v", err)
	}
	if l := len(comic.Comic); l != 1 {
		t.Errorf("Expect len of comic.Comic 1, got %d", l)
	}

	if v, ok := comic.Comic["test1"]; !ok || v != "1" {
		t.Errorf("Expect test1 in comic.Comic with value 1, got content %+v", comic.Comic)
	}
}

func TestToJSON(t *testing.T) {
	var comic Comics
	data, err := comic.ToJSON()
	if err != nil {
		t.Errorf("Expecting no error, got error: %+v", err)
	}
	exp := `{
  "comic": null,
  "anime": null
}`
	if s := string(data); s != exp {
		t.Errorf("Expect \"%s\", got \"%s\"", exp, s)
	}
}
