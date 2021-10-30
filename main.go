package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/hellflame/argparse"
	"hashkitty/algos"
	_ "hashkitty/rules"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type Settings struct {
	leftlist   *string
	wordlist   *string
	rules      *string
	potfile    *string
	attackMode *int
	hashType   *int

	tasks         *chan Task
	results       *chan Task
	potfileCloser *chan bool

	progress *sync.WaitGroup
	writes   *sync.WaitGroup
}

func NewSettings() *Settings {
	conf := argparse.ParserConfig{
		Usage:                  "",
		EpiLog:                 "",
		DisableHelp:            false,
		ContinueOnHelp:         false,
		DisableDefaultShowHelp: false,
		DefaultAction:          nil,
		AddShellCompletion:     false,
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
		Help: "Attack Mode",
	})

	hashType := parser.Int("m", "hash-type", &argparse.Option{
		Help: "Hash Type",
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
	return &Settings{leftlist, wordlist, rules, potfile, attackMode, hashType, &tasksChan, &goodTasksChan, &potfileCloser, &progress, &writes}
}

type Task struct {
	hash string
	salt string
	word string
}

func Worker(settings *Settings) {
	validator := algos.HASHCATALGOS[uint(*settings.hashType)]

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()
	for {
		select {
		case task := <-*settings.tasks:
			settings.progress.Done()
			if validator(task.hash, task.word, task.salt) {
				fmt.Printf("OK %s %s\n", task.hash, task.word)
				settings.writes.Add(1)
				*settings.results <- task
			}
		case <-ctx.Done():
			return
		}
	}
}

func PotfileWriter(settings *Settings) {
	potfile, err := os.OpenFile(*settings.potfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(777))
	if err != nil {
		panic(err)
	}
	for {
		select {
		case task := <-*settings.results:
			_, err := potfile.WriteString(fmt.Sprintf("%s:%s\n", task.hash, task.word))
			settings.writes.Done()
			if err != nil {
				panic(err)
			}
		case <-*settings.potfileCloser:
			potfile.Close()
			*settings.potfileCloser <- true
		}
	}
}
func PotfileCloser(settings *Settings) {
	*settings.potfileCloser <- true
	<-*settings.potfileCloser
}

func spawnWorkers(settings *Settings) {
	for i := runtime.NumCPU() * 10; i != 0; i-- {
		go Worker(settings)
	}
}

func main() {
	settings := NewSettings()

	spawnWorkers(settings)
	go PotfileWriter(settings)

	leftlist := NewLeftlist(settings)
	defer leftlist.Close()
	wordlist := NewWordlist(settings)
	defer wordlist.Close()
	ruleset := NewRuleset(settings)
	defer ruleset.Close()

	if *settings.attackMode == 0 {
		fmt.Println("Start attack mode 0")
		mode0(settings, leftlist, wordlist, ruleset)
	} else if *settings.attackMode == 9 {
		fmt.Println("Start attack mode 9")
		mode9(settings, leftlist, wordlist, ruleset)
	}

	settings.progress.Wait()
	settings.writes.Wait()
	PotfileCloser(settings)
}
