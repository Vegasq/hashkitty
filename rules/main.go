/*
Package rules

Library that implements subset of rules from hashcat.
*/
package rules

import (
	"strings"
)

// https://hashcat.net/wiki/doku.php?id=rule_based_attack

var simpleRulesMap = map[string]func(string, *strings.Reader) string{
	":": pass,
	"l": lowercase,
	"u": uppercase,
	"c": capitalize,
	"C": invertCapitalize,
	"t": toggleCase,
	"d": duplicate,
	"r": reverse,
	"f": reflect,
	"{": rotateLeft,
	"}": rotateRight,
	"[": truncateLeft,
	"]": truncateRight,
	"q": duplicateAll,
	"T": toggleAt,
	"p": duplicateTimes,
	"$": appendChar,
	"^": prependChar,
	"D": deleteAt,
	"x": extractRange,
	"O": omitRange,
	"i": insertAt,
	"o": overwriteAt,
	"'": truncateAt,
	"s": replace,
	"@": purge,
	"z": duplicateFirst,
	"Z": duplicateLast,
}

//Apply - applies rules to the given string
// Example:
// 		rules.Apply("abc", "$1") == "abc1"
func Apply(rule string, word string) string {
	ruleReader := strings.NewReader(rule)
	for {
		bt, err := ruleReader.ReadByte()
		if err != nil {
			break
		}
		modifier := simpleRulesMap[string(bt)]
		if modifier != nil {
			word = modifier(word, ruleReader)
		}
	}
	return word
}
