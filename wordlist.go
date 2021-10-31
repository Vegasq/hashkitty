package main

import (
	"bufio"
	"os"
)

type Wordlist struct {
	TextFileProcessor
}

func NewWordlist(settings *Settings) *Wordlist {
	fl, err := os.Open(*settings.wordlist)
	reader := bufio.NewReader(fl)
	if err != nil {
		panic("Failed to open wordlist")
	}
	return &Wordlist{TextFileProcessor{settings.wordlist, fl, reader}}
}
