package main

import (
	"context"
	"github.com/vegasq/hashkitty/algos"
	"log"
	"runtime"
	"time"
)

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
				cracked[sliceToArray(task.hash)] = true

				settings.writes.Add(1)
				*settings.results <- task
			}
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
