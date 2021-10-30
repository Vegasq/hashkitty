package rules

import (
	"strconv"
	"strings"
	"unicode"
)

func readNum(reader *strings.Reader) uint {
	r, _, err := reader.ReadRune()
	if err != nil {
		panic("Failed to read rule")
	}
	val := ruleNumericSystemToUint(r)
	return val
}

func readChar(reader *strings.Reader) rune {
	var char rune
	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		return 0
	}

	r, _, err := reader.ReadRune()
	if err != nil {
		return r
	}
	// Edge case to be fixed: $\
	if r == '\\' {
		// read x
		reader.ReadRune()
		h1, err := reader.ReadByte()
		if err != nil {
			return r
		}
		h2, err := reader.ReadByte()
		if err != nil {
			return r
		}
		char = rune(hexToByte(h1)<<4 + hexToByte(h2))
	} else {
		char = r
	}

	return char
}

func ruleNumericSystemToUint(r rune) uint {
	if unicode.IsNumber(r) {
		val, _ := strconv.Atoi(string(r))
		return uint(val)
	} else if unicode.IsLetter(r) {
		if unicode.IsLower(r) {
			r = unicode.ToUpper(r)
		}
		// From official documentation
		// https://hashcat.net/wiki/doku.php?id=rule_based_attack
		//  * Indicates that N starts at 0. For character positions other than 0-9 use A-Z (A=10)
		r = r - 55
	}
	return uint(r)
}
