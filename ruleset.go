package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"
)

type Ruleset struct {
	name   *string
	fl     *os.File
	reader *bufio.Reader
}

func (r *Ruleset) Reset() {
	r.fl.Seek(0, 0)
}

func (r *Ruleset) Close() {
	err := r.fl.Close()
	if err != nil && r.name != nil {
		log.Printf("Failed to close file %s", *r.name)
	}
}

func (r *Ruleset) GetNextRule() (string, error) {
	var rule, err = r.reader.ReadString('\n')
	rule = strings.TrimRight(rule, "\r\n")
	if err != nil {
		return rule, errors.New("failed to read a rule")
	}

	return rule, nil
}

func NewRuleset(settings *Settings) *Ruleset {
	if len(*settings.rules) > 0 {
		fl, err := os.Open(*settings.rules)
		rsReader := bufio.NewReader(fl)
		if err != nil {
			panic("Failed to open ruleset")
		}
		return &Ruleset{settings.rules, fl, rsReader}
	}
	return &Ruleset{nil, nil, nil}
}
