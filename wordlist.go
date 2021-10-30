package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Wordlist struct {
	name   *string
	fl     *os.File
	reader *bufio.Reader
}

func (w *Wordlist) Close() {
	err := w.fl.Close()
	if err != nil {
		fmt.Printf("Failed to close file %s", *w.name)
	}
}

func (w *Wordlist) Reset() {
	w.fl.Seek(0, 0)
}

func (w *Wordlist) GetNextLine() (string, error) {
	if w.reader == nil {
		w.reader = bufio.NewReader(w.fl)
	}

	var word, err = w.reader.ReadString('\n')
	word = strings.TrimRight(word, "\r\n")
	if err != nil {
		return word, errors.New("failed to read a word")
	}

	return word, nil
}

func NewWordlist(settings *Settings) *Wordlist {
	fl, err := os.Open(*settings.wordlist)
	reader := bufio.NewReader(fl)
	if err != nil {
		panic("Failed to open wordlist")
	}
	return &Wordlist{settings.wordlist, fl, reader}
}
