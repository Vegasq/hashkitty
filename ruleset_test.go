package main

import (
	"fmt"
	"os"
	"testing"
)

func TestNewRuleset(t *testing.T) {
	rulesetName := "mockruleset.txt"

	line1 := "$1 $2 $3"
	line2 := "{ { {"
	line3 := "] ] ]"

	fl, err := os.Create(rulesetName)
	if err != nil {
		panic(err)
	}
	fl.WriteString(fmt.Sprintf("%s\n%s\n%s", line1, line2, line3))
	fl.Close()

	defer os.Remove(rulesetName)

	settings := Settings{rules: &rulesetName}
	rs := NewRuleset(&settings)
	rs.Reset()

	for _, l := range []string{line1, line2, line3} {
		rule, _ := rs.GetNextRule()
		if l != rule {
			t.Errorf("Incorrect line in rules reader %s != %s", rule, l)
		}
	}
	rs.Close()
}
