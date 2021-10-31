package main

import (
	"log"
)

type Task struct {
	hash string
	salt string
	word string
}

func main() {
	settings := NewSettings()

	spawnWorkers(settings)
	go potfileWriter(settings)

	leftlist := NewLeftlist(settings)
	wordlist := NewWordlist(settings)
	defer wordlist.Close()
	ruleset := NewRuleset(settings)
	defer ruleset.Close()

	log.Printf("Start attack mode %d\n", *settings.attackMode)
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
