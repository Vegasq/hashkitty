package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type Leftlist struct {
	name   *string
	fl     *os.File
	reader *bufio.Reader
}

type LeftlistRecord struct {
	hash string
	salt string
}

func (l *Leftlist) Close() {
	err := l.fl.Close()
	if err != nil {
		log.Printf("Failed to close file %s", *l.name)
	}
}

func (l *Leftlist) GetNextRecord() (LeftlistRecord, error) {
	var line, eof = l.reader.ReadString('\n')
	line = strings.TrimRight(line, "\r\n")

	var hash, salt string
	if len(line) > 0 {
		dividedLine := strings.SplitN(line, ":", 2)
		hash = dividedLine[0]
		if len(dividedLine) == 2 {
			salt = dividedLine[1]
		}
	}
	return LeftlistRecord{hash, salt}, eof
}

func NewLeftlist(settings *Settings) *Leftlist {
	fl, err := os.Open(*settings.leftlist)
	if err != nil {
		panic("Failed to open leftlist")
	}
	return &Leftlist{settings.leftlist, fl, bufio.NewReader(fl)}
}
