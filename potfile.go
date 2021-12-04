package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func stringToHex(s string) string {
	var result string
	for _, c := range s {
		result += fmt.Sprintf("%02x", c)
	}
	return result
}

func taskToHashSaltRepr(settings *Settings, task Task) string {
	if len(task.salt) > 0 {
		salt := task.salt
		if *settings.hexSalt {
			salt = stringToHex(task.salt)
		}
		return fmt.Sprintf("%s:%s", task.hash, salt)
	} else {
		return fmt.Sprintf("%s", task.hash)
	}
}

func generateOutfileLine(settings *Settings, task Task) string {
	va := (*OutFileFormat)
	parts := strings.Split(va, ",")

	components := []string{}
	for i := 0; i < len(parts); i++ {
		val, err := strconv.Atoi(parts[i])
		if err != nil {
			continue
		}
		if val == 1 {
			components = append(components, taskToHashSaltRepr(settings, task))
		} else if val == 2 {
			components = append(components, task.word)
		} else if val == 3 {
			components = append(components, stringToHex(task.word))
		}
	}

	return strings.Join(components, ":")
}

func writeToPotfile(settings *Settings, potfile *os.File, task Task) error {
	var err error
	line := generateOutfileLine(settings, task)
	_, err = potfile.WriteString(fmt.Sprintf("%s\n", line))
	return err
}

func potfileWriter(settings *Settings) {
	if len(*OutFile) > 0 {
		*settings.potfile = *OutFile
	}

	potfile, err := os.OpenFile(*settings.potfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(777))
	if err != nil {
		panic(err)
	}
	for {
		select {
		case task := <-*settings.results:
			err := writeToPotfile(settings, potfile, task)
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
