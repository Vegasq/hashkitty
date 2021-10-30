package main

import (
	"hashkitty/rules"
	"log"
)


func sendTask(settings *Settings, hash LeftlistRecord, word string) {
	log.Println("CHECK " + hash.hash + " " + hash.salt + " " + word)

	cracked := *settings.cracked
	if cracked[hash.hash[0:32]] == true {
		log.Println("Skipping already known hash")
		return
	}

	settings.progress.Add(1)
	*(settings.tasks) <- Task{
		hash: hash.hash,
		salt: hash.salt,
		word: word,
	}
}

func combineWordWithRules(settings *Settings, ruleset *Ruleset, word string, hash LeftlistRecord) {
	ruleset.Reset()
	for {
		var rule, err = ruleset.GetNextRule()

		if len(rule) > 0 {
			processedWord := rules.Apply(rule, word)
			sendTask(settings, hash, processedWord)
		}

		if err != nil {
			break
		}
	}
}

func combineHashWithWords(settings *Settings, hash LeftlistRecord, wordlist *Wordlist, ruleset *Ruleset) {
	wordlist.Reset()
	for {
		word, eof := wordlist.GetNextLine()

		if len(word) > 0 {
			if ruleset.name != nil {
				combineWordWithRules(settings, ruleset, word, hash)
			} else {
				sendTask(settings, hash, word)
			}
		}

		if eof != nil {
			break
		}
	}
}

func readLeftlist(settings *Settings, leftlist *Leftlist, wordlist *Wordlist, ruleset *Ruleset) {
	for {
		hash, eof := leftlist.GetNextRecord()
		if len(hash.hash) > 0 {
			combineHashWithWords(settings, hash, wordlist, ruleset)
		}
		if eof != nil {
			return
		}
	}
}

func mode0(settings *Settings, leftlist *Leftlist, wordlist *Wordlist, ruleset *Ruleset) {
	readLeftlist(settings, leftlist, wordlist, ruleset)
}
