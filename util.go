package dbshiftcore

import (
	"fmt"
)

const successCharacter rune = '✔'
const failureCharacter rune = '✘'

// printSuccess prints a formatted text adding a special success character
func printSuccess(text string, args ...interface{}) {
	if len(args) > 0 {
		text = fmt.Sprintf(text, args...)
	}
	fmt.Printf("%c %s\n", successCharacter, text)
}

// printFailure prints a formatted text adding a special failure character
func printFailure(text string, args ...interface{}) {
	if len(args) > 0 {
		text = fmt.Sprintf(text, args...)
	}
	fmt.Printf("%c %s\n", failureCharacter, text)
}
