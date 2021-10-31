package main

func mode9(settings *Settings, leftlist *Leftlist, wordlist *Wordlist, ruleset *Ruleset) {
	for {
		hash, llEOF := leftlist.GetNextRecord()
		word, wlEOF := wordlist.ReadString()

		if ruleset.name != nil {
			combineWordWithRules(settings, ruleset, word, hash)
		} else {
			sendTask(settings, hash, word)
		}

		if llEOF != nil || wlEOF != nil {
			//fmt.Println("Exiting mode 9")
			return
		}
	}
}
