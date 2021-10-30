package rules

import "testing"

func TestApply(t *testing.T) {
	val := Apply("} } p2 r f", "p@ssW0rd")
	if val != "0Wss@pdr0Wss@pdr0Wss@pdrrdp@ssW0rdp@ssW0rdp@ssW0" {
		t.Errorf(val)
	}
}
