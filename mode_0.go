package main

import (
	"bufio"
	"hashkitty/rules"
	"strings"
)

func combineWordWithRules(settings *Settings, ruleset *Ruleset, word string, hash LeftlistRecord) {
	rsReader := bufio.NewReader(ruleset.fl)
	ruleset.fl.Seek(0, 0)
	for {
		var rule, err = rsReader.ReadString('\n')

		if len(rule) > 0 {
			//fmt.Println("Process rule", rule)
			rule = strings.Replace(rule, "\n", "", 1)
			rule = strings.Replace(rule, "\r", "", 1)
			processedWord := rules.Apply(rule, word)

			sendTask(settings, hash, processedWord)
		}

		if err != nil {
			break
		}
	}
}

func combineHashWithWords(settings *Settings, hash LeftlistRecord, wordlist *Wordlist, ruleset *Ruleset) {
	wlReader := bufio.NewReader(wordlist.fl)
	wordlist.fl.Seek(0, 0)
	for {
		var word, err = wlReader.ReadString('\n')
		word = strings.Replace(word, "\n", "", -1)
		word = strings.Replace(word, "\r", "", -1)

		if len(word) > 0 {
			if ruleset.name != nil {
				combineWordWithRules(settings, ruleset, word, hash)
			} else {
				sendTask(settings, hash, word)
			}
		}

		if err != nil {
			break
		}
	}
}

func readLeftlist(settings *Settings, leftlist *Leftlist, wordlist *Wordlist, ruleset *Ruleset) {
	for {
		hash, err := leftlist.GetNextRecord()
		if err != nil {
			return
		}
		combineHashWithWords(settings, hash, wordlist, ruleset)
	}
}

func mode0(settings *Settings, leftlist *Leftlist, wordlist *Wordlist, ruleset *Ruleset) {
	readLeftlist(settings, leftlist, wordlist, ruleset)
}
