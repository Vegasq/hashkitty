package main

import (
	"hashkitty/rules"
	"log"
)

func sendTask(settings *Settings, hash LeftlistRecord, word string) {
	settings.progress.Add(1)
	log.Println("CHECK " + hash.hash + " " + hash.salt + " " + word)

	*(settings.tasks) <- Task{
		hash: hash.hash,
		salt: hash.salt,
		word: word,
	}
}

func sendRulifiedTasks(settings *Settings, hash LeftlistRecord, word string, ruleset *Ruleset) {
	ruleset.Reset()
	for {
		rule, err := ruleset.GetNextRule()
		if err != nil {
			return
		}
		subWord := rules.Apply(rule, word)
		sendTask(settings, hash, subWord)
	}
}

func mode9(settings *Settings, leftlist *Leftlist, wordlist *Wordlist, ruleset *Ruleset) {
	for {
		hash, llEOF := leftlist.GetNextRecord()
		word, wlEOF := wordlist.GetNextLine()

		if ruleset.name != nil {
			//fmt.Println("send rulified")
			sendRulifiedTasks(settings, hash, word, ruleset)
		} else {
			sendTask(settings, hash, word)
		}

		if llEOF != nil || wlEOF != nil {
			//fmt.Println("Exiting mode 9")
			return
		}
	}
}
