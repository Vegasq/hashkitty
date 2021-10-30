package main

import (
	"fmt"
	"os"
	"testing"
)

func TestNewLeftlist(t *testing.T) {
	leftlistName := "mockleftlist.txt"

	line1 := "d131dd02c5e6eec4"
	line2 := "55ad340609f4b302"
	line3 := "e99f33420f577ee8"

	fl, err := os.Create(leftlistName)
	if err != nil {
		panic(err)
	}
	fl.WriteString(fmt.Sprintf("%s\n%s\n%s", line1, line2, line3))
	fl.Close()

	defer os.Remove(leftlistName)

	settings := Settings{leftlist: &leftlistName}
	ll := NewLeftlist(&settings)

	for _, l := range []string{line1, line2, line3, ""} {
		record, _ := ll.GetNextRecord()
		if l != record.hash {
			t.Errorf("Incorrect line in leftlist reader %s != %s", record.hash, l)
		}
	}
	ll.Close()
}
