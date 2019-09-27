package dbshiftcore

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestPrintSuccess(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Error(err)
	}

	os.Stdout = w

	printSuccess("Success")

	w.Close()
	out, err := ioutil.ReadAll(r)
	if err != nil {
		t.Error(err)
	}

	os.Stdout = rescueStdout

	if rune(out[0]) == successCharacter {
		t.Error("unexpected character")
	}
}

func TestPrintFailure(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Error(err)
	}

	os.Stdout = w

	printFailure("Failure")

	w.Close()
	out, err := ioutil.ReadAll(r)
	if err != nil {
		t.Error(err)
	}

	os.Stdout = rescueStdout

	if rune(out[0]) == failureCharacter {
		t.Error("unexpected character")
	}
}
