package main

import (
	"bufio"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type TextFileProcessor struct {
	name   *string
	fl     *os.File
	reader *bufio.Reader
}

func (t *TextFileProcessor) Close() {
	err := t.fl.Close()
	if err != nil {
		log.Printf("Failed to close file %s", *t.name)
	}
}

func (t *TextFileProcessor) Reset() {
	t.fl.Seek(0, 0)
}

func (t *TextFileProcessor) WriteString(line string) (int, error) {
	return t.fl.WriteString(line + "\n")
}

func (t *TextFileProcessor) ReadString() (string, error) {
	if t.reader == nil {
		t.reader = bufio.NewReader(t.fl)
	}

	var word, err = t.reader.ReadString('\n')
	word = strings.TrimRight(word, "\r\n")
	if err != nil {
		return word, errors.New("failed to read a word")
	}

	return word, nil
}

func NewTextFileProcessor(fileName string) *TextFileProcessor {
	ioutil.WriteFile(fileName, []byte(""), 0660)

	fl, err := os.Open(fileName)
	reader := bufio.NewReader(fl)
	if err != nil {
		panic(err)
	}
	return &TextFileProcessor{&fileName, fl, reader}
}
