package main

import (
	"errors"
	"fmt"
	"github.com/hellflame/argparse"
	"github.com/vegasq/hashkitty/algos"
	"os"
	"path/filepath"
	"sync"
)

func isValidAlgo(hashType uint) bool {
	val := algos.HASHCATALGOS[hashType]
	return val != nil
}

type Settings struct {
	leftlist   *string
	wordlist   *string
	rules      *string
	potfile    *string
	attackMode *int
	hashType   *int
	hexSalt    *bool
	remove     *bool

	tasks         *chan Task
	results       *chan Task
	potfileCloser *chan bool

	progress *sync.WaitGroup
	writes   *sync.WaitGroup

	crackedMap *sync.Map
	checked    *uint32
	maxGuesses uint32
}

func NewSettings() (*Settings, error) {
	conf := argparse.ParserConfig{
		Usage:                  "",
		EpiLog:                 "",
		DisableHelp:            false,
		ContinueOnHelp:         false,
		DisableDefaultShowHelp: false,
		DefaultAction:          nil,
		AddShellCompletion:     true,
		WithHint:               false,
	}
	parser := argparse.NewParser("HashKitty", "Hash cracking tool", &conf)

	parser.String("", "hashkitty", &argparse.Option{
		Positional: true,
		HideEntry:  true,
	})

	leftlist := parser.String("", "leftlist", &argparse.Option{
		Positional: true,
		Required:   true,
		Help:       "Leftlist file location",
	})

	wordlist := parser.String("", "wordlist", &argparse.Option{
		Positional: true,
		Required:   true,
		Help:       "Wordlist file location",
	})

	rules := parser.String("r", "rules-file", &argparse.Option{
		Help: "Rules file location",
	})

	potfile := parser.String("p", "potfile-path", &argparse.Option{
		Help:    "Potfile location",
		Default: "potfile.txt",
	})

	attackMode := parser.Int("a", "attack-mode", &argparse.Option{
		Help:     "Attack Mode",
		Required: true,
	})

	hashType := parser.Int("m", "hash-type", &argparse.Option{
		Help:     "Hash Type",
		Required: true,
	})

	remove := parser.Flag("", "remove", &argparse.Option{
		Help: "Enable removal of hashes once they are cracked",
	})

	hexSalt := parser.Flag("", "hex-salt", &argparse.Option{
		Help: "Salts provided in hex",
	})

	err := parser.Parse(os.Args)
	if err != nil {
		parser.PrintHelp()
		return &Settings{}, err
	}

	if isValidAlgo(uint(*hashType)) == false {
		return &Settings{}, errors.New(fmt.Sprintf("unknown mode %d", *hashType))
	}

	if _, err := os.Stat(*potfile); errors.Is(err, os.ErrNotExist) {
		os.MkdirAll(filepath.Dir(*potfile), 0660)
		fl, err := os.Create(*potfile)
		if err != nil {
			panic(err)
		}
		fl.Close()
	}

	progress := sync.WaitGroup{}

	writes := sync.WaitGroup{}
	tasksChan := make(chan Task)
	goodTasksChan := make(chan Task)
	potfileCloser := make(chan bool)

	var checked uint32 = 0
	var maxGuesses uint32 = 0

	s := &Settings{
		leftlist,
		wordlist,
		rules,
		potfile,
		attackMode,
		hashType,
		hexSalt,
		remove,
		&tasksChan,
		&goodTasksChan,
		&potfileCloser,
		&progress,
		&writes,
		&sync.Map{},
		&checked,
		maxGuesses,
	}
	s.maxGuesses = calculateMaxGuesses(s)
	return s, nil
}
