package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/hellflame/argparse"
	"hashkitty/algos"
	_ "hashkitty/rules"
	"os"
	"runtime"
	"strings"
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

	tasks   *chan Task
	results *chan Task

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
		Help: "Potfile location",
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
	progress := sync.WaitGroup{}
	writes := sync.WaitGroup{}
	tasksChan := make(chan Task)
	goodTasksChan := make(chan Task)
	return &Settings{leftlist, wordlist, rules, potfile, attackMode, hashType, &tasksChan, &goodTasksChan, &progress, &writes}
}

type Leftlist struct {
	name *string
	fl   *os.File

	// attack mode 9
	reader *bufio.Reader
}

type LeftlistRecord struct {
	hash string
	salt string
}

func (l *Leftlist) GetNextRecord() (LeftlistRecord, error) {
	if l.reader == nil {
		l.reader = bufio.NewReader(l.fl)
	}

	var hash, eof = l.reader.ReadString('\n')
	if len(hash) > 0 {
		hash = strings.TrimRight(hash, "\r\n")
		dividedHash := strings.SplitN(hash, ":", 2)
		if len(dividedHash) == 2 {
			return LeftlistRecord{dividedHash[0], dividedHash[1]}, nil
		}
		return LeftlistRecord{hash, ""}, nil
	}
	return LeftlistRecord{}, eof
}

func NewLeftlist(settings *Settings) *Leftlist {
	//fmt.Println("leftlist", *settings.leftlist)
	fl, err := os.Open(*settings.leftlist)
	if err != nil {
		panic("Failed to open leftlist")
	}
	return &Leftlist{settings.leftlist, fl, nil}
}

type Ruleset struct {
	name *string
	fl   *os.File

	// attack mode 9
	reader *bufio.Reader
}

func (r *Ruleset) Reset() {
	r.fl.Seek(0, 0)
}
func (r *Ruleset) GetNextRule() (string, error) {
	if r.reader == nil {
		r.reader = bufio.NewReader(r.fl)
	}

	var rule, err = r.reader.ReadString('\n')
	if err != nil {
		return "", errors.New("failed to read a rule")
	}
	rule = strings.TrimRight(rule, "\r\n")

	return rule, nil
}

func NewRuleset(settings *Settings) *Ruleset {
	if len(*settings.rules) > 0 {
		fl, err := os.Open(*settings.rules)
		if err != nil {
			panic("Failed to open ruleset")
		}
		return &Ruleset{settings.rules, fl, nil}
	}
	return &Ruleset{nil, nil, nil}
}

type Wordlist struct {
	name *string
	fl   *os.File

	// attack mode 9
	reader *bufio.Reader
}

func (w *Wordlist) GetNextLine() (string, error) {
	if w.reader == nil {
		w.reader = bufio.NewReader(w.fl)
	}

	var word, err = w.reader.ReadString('\n')
	word = strings.TrimRight(word, "\r\n")
	if err != nil {
		return "", errors.New("failed to read a word")
	}

	return word, nil
}

func NewWordlist(settings *Settings) *Wordlist {
	//fmt.Println("wordlist", *settings.wordlist)
	fl, err := os.Open(*settings.wordlist)
	if err != nil {
		panic("Failed to open wordlist")
	}
	return &Wordlist{settings.wordlist, fl, nil}
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
	defer potfile.Close()
	for {
		task := <-*settings.results
		//fmt.Println(task.word, task.hash)
		n, err := potfile.WriteString(fmt.Sprintf("%s:%s\n", task.hash, task.word))
		settings.writes.Done()
		if err != nil {
			fmt.Println(n)
			panic(err)
		}
	}
}

func main() {
	settings := NewSettings()
	for i := runtime.NumCPU() * 10; i != 0; i-- {
		go Worker(settings)
	}

	go PotfileWriter(settings)

	leftlist := NewLeftlist(settings)
	defer leftlist.fl.Close()
	wordlist := NewWordlist(settings)
	defer wordlist.fl.Close()
	ruleset := NewRuleset(settings)
	defer ruleset.fl.Close()

	if *settings.attackMode == 0 {
		fmt.Println("Start attack mode 0")
		mode0(settings, leftlist, wordlist, ruleset)
	} else if *settings.attackMode == 9 {
		fmt.Println("Start attack mode 9")
		mode9(settings, leftlist, wordlist, ruleset)
	}

	settings.progress.Wait()
	settings.writes.Wait()
}
