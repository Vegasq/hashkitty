/*
Package main

Reimplementation of _some_ of the [HashCat](https://github.com/hashcat/hashcat) features in GO.
*/
package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/debug"
	"time"
)

type Task struct {
	hash string
	salt string
	word string
}

func (t *Task) toString() string {
	return fmt.Sprintf("%s:%s:%s", t.hash, t.salt, t.word)
}

//https://stackoverflow.com/questions/24562942/golang-how-do-i-determine-the-number-of-lines-in-a-file-efficiently/24563853
func lineCounter(r io.Reader) (uint32, error) {
	buf := make([]byte, 32*1024)
	count := uint32(0)
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += uint32(bytes.Count(buf[:c], lineSep))

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func calculateMaxGuesses(settings *Settings) uint32 {
	ll, _ := os.Open(*settings.leftlist)
	llLines, _ := lineCounter(ll)
	ll.Close()
	llLines++

	fmt.Println(*settings.wordlist)
	wl, _ := os.Open(*settings.wordlist)
	wlLines, _ := lineCounter(wl)
	wl.Close()
	wlLines++

	rsLines := uint32(1)
	if *settings.rules != "" {
		rs, _ := os.Open(*settings.rules)
		rsLines, _ = lineCounter(rs)
		rs.Close()
	}

	return llLines * wlLines * rsLines
}

func GC() {
	t := time.Now()
	for {
		if time.Since(t) > time.Second*30 {
			t = time.Now()
			log.Println("Garbage collector")
			debug.FreeOSMemory()
			log.Println("Garbage collector done")
		}
	}
}

func main() {
	go GC()

	settings, err := NewSettings()
	if err != nil {
		fmt.Println(err)
		return
	}

	spawnWorkers(settings)
	go potfileWriter(settings)
	go checkedReporter(settings)

	leftlist := NewLeftlist(settings)
	wordlist := NewWordlist(settings)
	defer wordlist.Close()

	ruleset := NewRuleset(settings)
	defer ruleset.Close()

	log.Printf(
		"Cracking leftlist %s with wordlist %s and ruleset %s using mode %d\n",
		*settings.leftlist,
		*settings.wordlist,
		*settings.rules,
		*settings.attackMode,
	)
	if *settings.attackMode == 0 {
		mode0(settings, leftlist, wordlist, ruleset)
	} else if *settings.attackMode == 9 {
		mode9(settings, leftlist, wordlist, ruleset)
	}

	log.Println("Waiting for workers to finish")
	settings.progress.Wait()
	// Once settings.progress is done, we can close leftlist
	leftlist.Close()

	settings.writes.Wait()
	potfileCloser(settings)
	if *settings.remove {
		leftlist.CleanLeftlistWithPotfile(settings)
	}
}
