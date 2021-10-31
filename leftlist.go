package main

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Leftlist struct {
	TextFileProcessor

	hexSalt bool
}

type LeftlistRecord struct {
	hash string
	salt string
}

func isLineInFile(line string, potfile *TextFileProcessor) bool {
	potfile.Reset()
	for {
		ptLine, ptErr := potfile.ReadString()
		if ptLine == line {
			return true
		}
		if errors.Is(ptErr, io.EOF) {
			break
		}
	}
	return false
}

func (l *Leftlist) CleanLeftlistWithPotfile(settings *Settings) {
	potfile := NewTextFileProcessor(*settings.potfile)
	tmpLeftlist := NewTextFileProcessor(*settings.leftlist + "_tmp")
	leftlist := NewTextFileProcessor(*settings.leftlist)

	for {
		llLine, llErr := leftlist.ReadString()
		if isLineInFile(llLine, potfile) == false {
			_, nptErr := tmpLeftlist.WriteString(llLine)
			if nptErr != nil {
				log.Println(nptErr)
			}
		}
		if llErr != nil {
			break
		}
	}

	potfile.Close()
	leftlist.Close()
	tmpLeftlist.Close()

	copyFile(*settings.leftlist+"_tmp", *settings.leftlist)
}

func (l *Leftlist) GetNextRecord() (LeftlistRecord, error) {
	var line, eof = l.ReadString()

	var hash, salt string
	if len(line) > 0 {
		dividedLine := strings.SplitN(line, ":", 2)
		hash = dividedLine[0]
		if len(dividedLine) == 2 {
			salt = dividedLine[1]
		}
	}
	if l.hexSalt {
		salt = hexToString(salt)
	}
	return LeftlistRecord{hash, salt}, eof
}

func NewLeftlist(settings *Settings) *Leftlist {
	fl, err := os.Open(*settings.leftlist)
	if err != nil {
		panic("Failed to open leftlist")
	}
	return &Leftlist{TextFileProcessor{settings.leftlist, fl, bufio.NewReader(fl)}, *settings.hexSalt}
}

func hexToString(salt string) string {
	s, err := hex.DecodeString(salt)
	if err != nil {
		panic(err)
	}
	return string(s)
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
