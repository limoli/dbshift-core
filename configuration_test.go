package dbshiftcore

import "testing"

func TestGetConfiguration(t *testing.T) {

	if _, err := getConfiguration(); err == nil {
		t.Error("expected missing DBSHIFT_ABS_FOLDER_MIGRATIONS environment variable")
	}

	// TODO: More tests about configuration and params availability

}
