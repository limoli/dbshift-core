package dbshiftcore

import (
	"fmt"
)

const successCharacter rune = '✔'
const failureCharacter rune = '✘'

// PrintSuccess prints a formatted text adding a special success character
func PrintSuccess(text string, args ...interface{}) {
	if len(args) > 0 {
		text = fmt.Sprintf(text, args...)
	}
	fmt.Printf("%c %s\n", successCharacter, text)
}

// PrintFailure prints a formatted text adding a special failure character
func PrintFailure(text string, args ...interface{}) {
	if len(args) > 0 {
		text = fmt.Sprintf(text, args...)
	}
	fmt.Printf("%c %s\n", failureCharacter, text)
}
