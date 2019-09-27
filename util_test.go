package dbshiftcore

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestPrintSuccess(t *testing.T) {

	// Formatted string => args
	tests := map[string][]interface{}{
		"Success":         nil,
		"Success: %s":     {"flott!"},
		"Success: %d %s!": {0, "errors"},
	}

	expectedResult := map[string]string{
		"Success":         "✔ Success\n",
		"Success: %s":     "✔ Success: flott!\n",
		"Success: %d %s!": "✔ Success: 0 errors!\n",
	}

	for k, v := range tests {
		rescueStdout := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			t.Error(err)
		}

		os.Stdout = w

		printSuccess(k, v...)

		w.Close()
		out, err := ioutil.ReadAll(r)
		if err != nil {
			t.Error(err)
		}

		os.Stdout = rescueStdout

		output := string(out)
		runes := []rune(output)

		if runes[0] != successCharacter {
			t.Errorf("unexpected character %c instead of %c", runes[0], successCharacter)
		}

		if output != expectedResult[k] {
			t.Errorf("unexpected success message %s instead of %s", output, expectedResult[k])
		}
	}

}

func TestPrintFailure(t *testing.T) {

	// Formatted string => args
	tests := map[string][]interface{}{
		"Failure":         nil,
		"Failure: %s":     {"uffda!"},
		"Failure: %d %s!": {5, "errors"},
	}

	expectedResult := map[string]string{
		"Failure":         "✘ Failure\n",
		"Failure: %s":     "✘ Failure: uffda!\n",
		"Failure: %d %s!": "✘ Failure: 5 errors!\n",
	}

	for k, v := range tests {
		rescueStdout := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			t.Error(err)
		}

		os.Stdout = w

		printFailure(k, v...)

		w.Close()
		out, err := ioutil.ReadAll(r)
		if err != nil {
			t.Error(err)
		}

		os.Stdout = rescueStdout

		output := string(out)
		runes := []rune(output)

		if runes[0] != failureCharacter {
			t.Errorf("unexpected character %c instead of %c", runes[0], successCharacter)
		}

		if output != expectedResult[k] {
			t.Errorf("unexpected failure message %s instead of %s", output, expectedResult[k])
		}

	}
}
