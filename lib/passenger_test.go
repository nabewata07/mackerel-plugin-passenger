package mppassenger

import (
	"fmt"
	"testing"
)

func TestGenerateCmdAry(t *testing.T) {
	p := PassengerPlugin{}
	c := generateCmdAry(p)

	fmt.Printf("%v", c)

	if c[0] != "passenger-status" {
		t.Errorf("\nexpected command : %s\nactual : %s", "passenger-status", c[0])
	}
}

// func TestGenerateCmdAry(t *testing.T) {
// }
