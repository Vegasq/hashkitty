package main

import (
	"fmt"
	"log"
	"os"
)

func potfileWriter(settings *Settings) {
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
