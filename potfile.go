package main

import (
	"fmt"
	"log"
	"os"
)

// being read from status printer and changed by potfileWriter routines
// can be sometimes off, but we do not care
var savedRecordsCounter = uint32(0)
var recentlySaved = map[string]bool{}

func stringToHex(s string) string {
	var result string
	for _, c := range s {
		result += fmt.Sprintf("%02x", c)
	}
	return result
}

func writeToPotfile(settings *Settings, potfile *os.File, task Task) error {
	var err error
	if len(task.salt) > 0 {
		salt := task.salt
		if *settings.hexSalt {
			salt = stringToHex(task.salt)
		}
		_, err = potfile.WriteString(fmt.Sprintf("%s:%s:%s\n", task.hash, salt, task.word))
	} else {
		_, err = potfile.WriteString(fmt.Sprintf("%s:%s\n", task.hash, task.word))
	}
	return err
}

func isAlreadySaved(t *Task) bool {
	_, ok := recentlySaved[t.toString()]
	return ok
}

func markAsSaved(t *Task) {
	recentlySaved[t.toString()] = true
}

func potfileWriter(settings *Settings) {
	potfile, err := os.OpenFile(*settings.potfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.FileMode(777))
	if err != nil {
		panic(err)
	}
	for {
		select {
		case task := <-*settings.results:
			if isAlreadySaved(&task) {
				settings.writes.Done()
				continue
			}
			markAsSaved(&task)
			err := writeToPotfile(settings, potfile, task)
			savedRecordsCounter += 1
			settings.writes.Done()
			if err != nil {
				panic(err)
			}
		case <-*settings.potfileCloser:
			if err := potfile.Close(); err != nil {
				log.Printf("Failed to close potfile: %e\n", err)
			}
			*settings.potfileCloser <- true
		}
	}
}

func potfileCloser(settings *Settings) {
	*settings.potfileCloser <- true
	<-*settings.potfileCloser
	log.Printf("Potfile %s\n", *settings.potfile)
}
