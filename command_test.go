package dbshift_core

import (
	"testing"
)

func TestNewCmd(t *testing.T) {

	if _, err := NewCmd(nil); err == nil {
		t.Error("expected missing db implementation error")
	}

}
