package main

import (
	"errors"
	"github.com/hellflame/argparse"
	"os"
	"path/filepath"
	"sync"
)

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

	cracked      *map[[32]int32]bool
	crackedMutex *sync.RWMutex
}

func NewSettings() *Settings {
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
	parser.String("h", "huh", &argparse.Option{
		Positional: true,
		Required:   true,
		Help:       "huh",
	})
	leftlist := parser.String("l", "leftlist", &argparse.Option{
		Positional: true,
		Required:   true,
		Help:       "Leftlist file location",
	})

	wordlist := parser.String("w", "wordlist", &argparse.Option{
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

	hexSalt := parser.Flag("hs", "hex-salt", &argparse.Option{
		Help: "Salts provided in hex",
	})

	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		parser.PrintHelp()
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

	// Possible collision
	cracked := map[[32]int32]bool{}
	crackedMutex := sync.RWMutex{}

	return &Settings{
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
		&cracked,
		&crackedMutex,
	}
}
