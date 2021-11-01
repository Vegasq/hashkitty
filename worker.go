package main

import (
	"context"
	"github.com/vegasq/hashkitty/algos"
	"log"
	"runtime"
	"time"
)

var lastChecked = uint32(1)

func checkedReporter(settings *Settings) {
	t := time.Now()
	for {
		if time.Since(t) > time.Second*5 {
			checked := *settings.checked

			perMinute := (checked - lastChecked) * 12
			if perMinute == 0 {
				perMinute = 1
			}
			left := (settings.maxGuesses - checked) / perMinute

			t = time.Now()
			log.Printf("Checked:\t%d/%d\tCracked:\t%d\tSpeed per minute:\t%d\tMinutes left:\t%d\tWorkers:\t%d\n",
				checked, settings.maxGuesses, savedRecordsCounter, perMinute, left, runtime.NumGoroutine())

			lastChecked = checked
		}
	}
}

func worker(settings *Settings) {
	validator := algos.HASHCATALGOS[uint(*settings.hashType)]

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()
	for {
		select {
		case task := <-*settings.tasks:
			k := sliceToArray(&task)
			_, isCracked := settings.crackedMap.Load(k)
			if isCracked == false && validator(task.hash, task.word, task.salt) {
				settings.crackedMap.Store(k, true)
				settings.writes.Add(1)
				*settings.results <- task
			}

			*settings.checked++
			settings.progress.Done()
		case <-ctx.Done():
			log.Println("Closing worker")
			return
		}
	}
}

func spawnWorkers(settings *Settings) {
	log.Printf("Routines before workers spawn: %d", runtime.NumGoroutine())
	for i := runtime.NumCPU() * 10; i != 0; i-- {
		go worker(settings)
	}
	log.Printf("Routines after workers spawn: %d", runtime.NumGoroutine())
}
