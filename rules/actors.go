package rules

import (
	"fmt"
	"strings"
	"unicode"
)

func pass(str string, reader *strings.Reader) string {
	return str
}

func lowercase(str string, reader *strings.Reader) string {
	return strings.ToLower(str)
}

func uppercase(str string, reader *strings.Reader) string {
	return strings.ToUpper(str)
}

func capitalize(str string, reader *strings.Reader) string {
	wordReader := strings.NewReader(str)
	r, _, err := wordReader.ReadRune()
	if err != nil {
		return str
	}
	word := strings.ToUpper(string(r))

	for {
		r, _, err := wordReader.ReadRune()
		if err != nil {
			break
		}
		word += strings.ToLower(string(r))
	}
	return word
}

func invertCapitalize(str string, reader *strings.Reader) string {
	wordReader := strings.NewReader(str)
	r, _, err := wordReader.ReadRune()
	if err != nil {
		fmt.Println("Failed to read rune in capitalize")
		return str
	}
	word := strings.ToLower(string(r))

	for {
		r, _, err := wordReader.ReadRune()
		if err != nil {
			break
		}
		word += strings.ToUpper(string(r))
	}
	return word
}

func toggleCase(str string, reader *strings.Reader) string {
	word := ""
	wordReader := strings.NewReader(str)
	for {
		r, _, err := wordReader.ReadRune()
		if err != nil {
			break
		}
		if unicode.IsLetter(r) {
			if unicode.IsLower(r) {
				word += string(unicode.ToUpper(r))
			}
			if unicode.IsUpper(r) {
				word += string(unicode.ToLower(r))
			}
		} else {
			word += string(r)
		}

	}
	return word
}

func toggleAt(str string, reader *strings.Reader) string {
	pos := readNum(reader)

	word := ""
	j := uint(0)
	wordReader := strings.NewReader(str)
	for {
		r, _, err := wordReader.ReadRune()
		if err != nil {
			break
		}

		if j == pos {
			word += string(unicode.ToUpper(r))
		} else {
			word += string(r)
		}
		j += 1
	}
	return word
}

func reverse(str string, reader *strings.Reader) string {
	word := ""
	for i := len(str) - 1; i >= 0; i-- {
		word += string(str[i])
	}
	return word
}

func duplicate(str string, reader *strings.Reader) string {
	return str + str
}

func duplicateTimes(str string, reader *strings.Reader) string {
	r, _, err := reader.ReadRune()
	if err != nil {
		panic("Failed to read rule")
	}
	val := ruleNumericSystemToUint(r)
	word := str
	for ; val > 0; val-- {
		word += str
	}
	return word
}

func reflect(str string, reader *strings.Reader) string {
	word := ""
	for i := len(str) - 1; i >= 0; i-- {
		word += string(str[i])
	}
	return str + word
}

func rotateLeft(str string, reader *strings.Reader) string {
	if len(str) == 0 {
		return ""
	}
	str += string(str[0])
	return str[1:]
}

func rotateRight(str string, reader *strings.Reader) string {
	if len(str) == 0 {
		return ""
	}
	last := string(str[len(str)-1])
	return last + str[:len(str)-1]
}

func appendChar(str string, reader *strings.Reader) string {
	r := readChar(reader)
	str = str + string(r)
	return str
}
func prependChar(str string, reader *strings.Reader) string {
	r := readChar(reader)
	str = string(r) + str
	return str
}

func truncateLeft(str string, reader *strings.Reader) string {
	if len(str) == 0 {
		return str
	}
	return str[1:]
}

func truncateRight(str string, reader *strings.Reader) string {
	if len(str) == 0 {
		return str
	}
	return str[:len(str)-1]
}

func deleteAt(str string, reader *strings.Reader) string {
	pos := readNum(reader)

	word := ""
	j := uint(0)
	wordReader := strings.NewReader(str)
	for {
		r, _, err := wordReader.ReadRune()
		if err != nil {
			break
		}

		if j != pos {
			word += string(r)
		}
		j += 1
	}
	return word
}

func extractRange(str string, reader *strings.Reader) string {
	startPos := readNum(reader)
	totalChars := readNum(reader)

	if len(str) < int(startPos) {
		return ""
	} else if len(str) < int(startPos+totalChars) {
		return str[startPos:]
	}
	return str[startPos : startPos+totalChars]
}

func omitRange(str string, reader *strings.Reader) string {
	startPos := readNum(reader)

	r, _, err := reader.ReadRune()
	if err != nil {
		panic("Failed to read rule")
	}
	totalChars := ruleNumericSystemToUint(r)
	if len(str) < int(startPos) {
		return str
	} else if len(str) < int(startPos+totalChars) {
		return str[:startPos]
	}

	return str[:startPos] + str[startPos+totalChars:]
}

func insertAt(str string, reader *strings.Reader) string {
	r, _, err := reader.ReadRune()
	if err != nil {
		panic("Failed to read rule")
	}
	pos := ruleNumericSystemToUint(r)

	r = readChar(reader)

	wordReader := strings.NewReader(str)
	i := uint(0)
	word := ""
	for {
		run, _, err := wordReader.ReadRune()
		if err != nil {
			break
		}
		if i == pos {
			word += string(r)
		}
		word += string(run)
		i++
	}
	return word
}

func overwriteAt(str string, reader *strings.Reader) string {
	pos := readNum(reader)
	r := readChar(reader)

	wordReader := strings.NewReader(str)
	i := uint(0)
	word := ""
	for {
		run, _, err := wordReader.ReadRune()
		if err != nil {
			break
		}
		if i == pos {
			word += string(r)
		} else {
			word += string(run)
		}
		i++
	}
	return word
}

func truncateAt(str string, reader *strings.Reader) string {
	r, _, err := reader.ReadRune()
	if err != nil {
		panic("Failed to read rule")
	}
	pos := ruleNumericSystemToUint(r)
	if len(str) <= int(pos) {
		return str
	}

	return str[:pos]
}

func replace(str string, reader *strings.Reader) string {
	repl := readChar(reader)
	with := readChar(reader)

	return strings.ReplaceAll(str, string(repl), string(with))
}

func purge(str string, reader *strings.Reader) string {
	repl := readChar(reader)

	return strings.ReplaceAll(str, string(repl), "")
}

func duplicateFirst(str string, reader *strings.Reader) string {
	count := readNum(reader)
	wordReader := strings.NewReader(str)
	word := ""

	for {
		r, _, err := wordReader.ReadRune()
		if err != nil {
			break
		}
		for ; count != 0; count-- {
			word += string(r)
		}
		word += string(r)
	}

	return word
}

func duplicateLast(str string, reader *strings.Reader) string {
	count := readNum(reader)
	wordReader := strings.NewReader(str)
	word := ""
	var lastRune rune

	for {
		r, _, err := wordReader.ReadRune()
		if err != nil {
			for ; count != 0; count-- {
				word += string(lastRune)
			}
			break
		}
		lastRune = r
		word += string(r)
	}

	return word
}

func duplicateAll(str string, reader *strings.Reader) string {
	wordReader := strings.NewReader(str)
	word := ""
	for {
		r, _, err := wordReader.ReadRune()
		if err != nil {
			break
		}
		word += string(r) + string(r)
	}
	return word
}
