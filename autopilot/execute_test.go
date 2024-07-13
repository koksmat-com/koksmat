package autopilot

import (
	"log"
	"testing"
)

func Test1(t *testing.T) {
	cb := func(isStdOut bool, output string) {
		if isStdOut {
			log.Println(output)
		} else {
			log.Println("Error:", output)
		}
	}

	_, err := Execute("ls", []string{"-la"}, Options{Timeout: 5, Cwd: ""}, cb)

	if err != nil {
		t.Error(err)
	} else {
		t.Log("done")
	}

}
