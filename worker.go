package main

import (
	"context"
	"github.com/vegasq/hashkitty/algos"
	"log"
	"runtime"
	"time"
)

func CheckedReporter(settings *Settings) {
	t := time.Now()
	for {
		if time.Since(t) > time.Second*5 {
			t = time.Now()
			log.Printf("Checked %d/%d", *settings.checked, settings.maxGuesses)
		}
	}
}

func Worker(settings *Settings) {
	validator := algos.HASHCATALGOS[uint(*settings.hashType)]

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()
	for {
		select {
		case task := <-*settings.tasks:
			if validator(task.hash, task.word, task.salt) {
				log.Printf("OK %s %s\n", task.hash, task.word)

				cracked := *settings.cracked
				settings.crackedMutex.Lock()
				cracked[sliceToArray(task.hash)] = true
				settings.crackedMutex.Unlock()

				settings.writes.Add(1)
				*settings.results <- task
			}

			settings.crackedMutex.Lock()
			*settings.checked++
			settings.crackedMutex.Unlock()

			settings.progress.Done()
		case <-ctx.Done():
			return
		}
	}
}

func spawnWorkers(settings *Settings) {
	for i := runtime.NumCPU() * 10; i != 0; i-- {
		go Worker(settings)
	}
}
