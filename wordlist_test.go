package main

import (
	"fmt"
	"os"
	"testing"
)

func TestNewWordlist(t *testing.T) {
	wordlistName := "mockwordlist.txt"

	line1 := "password1"
	line2 := "password2"
	line3 := "password3"

	fl, err := os.Create(wordlistName)
	if err != nil {
		panic(err)
	}
	fl.WriteString(fmt.Sprintf("%s\n%s\n%s", line1, line2, line3))
	fl.Close()

	defer os.Remove(wordlistName)

	settings := Settings{wordlist: &wordlistName}
	wl := NewWordlist(&settings)
	wl.Reset()

	for _, l := range []string{line1, line2, line3} {
		plain, _ := wl.GetNextLine()
		if l != plain {
			t.Errorf("Incorrect line in wordlist reader %s != %s", plain, l)
		}
	}
	wl.Close()
}
