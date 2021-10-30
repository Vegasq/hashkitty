package rules

import (
	reflectStd "reflect"
	"runtime"
	"strings"
	"testing"
)

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflectStd.ValueOf(i).Pointer()).Name()
}

func TestRulesActors(t *testing.T) {
	type Case struct {
		f      func(string, *strings.Reader) string
		input  string
		output string
		reader *strings.Reader
	}

	cases := []Case{
		{lowercase, "p@ssW0rd", "p@ssw0rd", strings.NewReader("")},
		{uppercase, "p@ssW0rd", "P@SSW0RD", strings.NewReader("")},
		{capitalize, "p@ssW0rd", "P@ssw0rd", strings.NewReader("")},
		{invertCapitalize, "p@ssW0rd", "p@SSW0RD", strings.NewReader("")},
		{toggleCase, "p@ssW0rd", "P@SSw0RD", strings.NewReader("")},
		{toggleAt, "p@ssW0rd", "p@sSW0rd", strings.NewReader("3")},
		{reverse, "p@ssW0rd", "dr0Wss@p", strings.NewReader("")},
		{duplicate, "p@ssW0rd", "p@ssW0rdp@ssW0rd", strings.NewReader("")},
		{duplicateTimes, "p@ssW0rd", "p@ssW0rdp@ssW0rdp@ssW0rd", strings.NewReader("2")},
		{reflect, "p@ssW0rd", "p@ssW0rddr0Wss@p", strings.NewReader("")},
		{rotateLeft, "p@ssW0rd", "@ssW0rdp", strings.NewReader("")},
		{rotateRight, "p@ssW0rd", "dp@ssW0r", strings.NewReader("")},
		{appendChar, "p@ssW0rd", "p@ssW0rd1", strings.NewReader("1")},
		{prependChar, "p@ssW0rd", "1p@ssW0rd", strings.NewReader("1")},
		{truncateLeft, "p@ssW0rd", "@ssW0rd", strings.NewReader("1")},
		{truncateRight, "p@ssW0rd", "p@ssW0r", strings.NewReader("1")},
		{deleteAt, "p@ssW0rd", "p@sW0rd", strings.NewReader("3")},
		{extractRange, "p@ssW0rd", "p@ss", strings.NewReader("04")},
		{omitRange, "p@ssW0rd", "psW0rd", strings.NewReader("12")},
		{insertAt, "p@ssW0rd", "p@ss!W0rd", strings.NewReader("4!")},
		{overwriteAt, "p@ssW0rd", "p@s$W0rd", strings.NewReader("3$")},
		{truncateAt, "p@ssW0rd", "p@ssW0", strings.NewReader("6")},
		{replace, "p@ssW0rd", "p@$$W0rd", strings.NewReader("s$")},
		{purge, "p@ssW0rd", "p@W0rd", strings.NewReader("s")},
		{duplicateFirst, "p@ssW0rd", "ppp@ssW0rd", strings.NewReader("2")},
		{duplicateLast, "p@ssW0rd", "p@ssW0rddd", strings.NewReader("2")},
		{duplicateAll, "p@ssW0rd", "pp@@ssssWW00rrdd", strings.NewReader("")},
	}

	for i := range cases {
		if got := cases[i].f(cases[i].input, cases[i].reader); got != cases[i].output {
			t.Errorf("%s = %v, want %v", GetFunctionName(cases[i].f), got, cases[i].output)
		}
	}
}
