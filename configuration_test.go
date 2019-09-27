package dbshiftcore

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetConfiguration(t *testing.T) {

	err := os.Unsetenv("DBSHIFT_ABS_FOLDER_MIGRATIONS")
	if err != nil {
		t.Error(err)
	}

	if _, err := getConfiguration(); err == nil {
		t.Error("expected missing DBSHIFT_ABS_FOLDER_MIGRATIONS environment variable")
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	migrationsPath := filepath.Join(wd, "example", "migrations")
	err = os.Setenv("DBSHIFT_ABS_FOLDER_MIGRATIONS", migrationsPath)
	if err != nil {
		t.Error(err)
	}

	cfg, err := getConfiguration()
	if err != nil {
		t.Error("expected set DBSHIFT_ABS_FOLDER_MIGRATIONS environment variable")
	} else if cfg.MigrationsPath != migrationsPath {
		t.Errorf("expected same migration path: %s != %s", cfg.MigrationsPath, migrationsPath)
	}

}

func TestGetConfigurationOption(t *testing.T) {

	wd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	migrationsPath := filepath.Join(wd, "example", "migrations")
	err = os.Setenv("DBSHIFT_ABS_FOLDER_MIGRATIONS", migrationsPath)
	if err != nil {
		t.Error(err)
	}

	err = os.Setenv("DBSHIFT_OPTION_IS_CREATE_DISABLED", "true")
	if err != nil {
		t.Error(err)
	}

	err = os.Setenv("DBSHIFT_OPTION_IS_DOWNGRADE_DISABLED", "true")
	if err != nil {
		t.Error(err)
	}

	err = os.Setenv("DBSHIFT_OPTION_IS_UPGRADE_DISABLED", "true")
	if err != nil {
		t.Error(err)
	}

	opts, err := getOptions()
	if err != nil {
		t.Error("expected set DBSHIFT_ABS_FOLDER_MIGRATIONS environment variable")
	} else if !opts.IsCreateDisabled {
		t.Error("expected true value for IsCreateDisabled")
	} else if !opts.IsDowngradeDisabled {
		t.Error("expected true value for IsDowngradeDisabled")
	} else if !opts.IsUpgradeDisabled {
		t.Error("expected true value for IsUpgradeDisabled")
	}
}

func TestGetEnvVar(t *testing.T) {
	if result, err := getEnvVar("unavailable_environment_variable!"); err == nil || result != "" {
		t.Error("expected missing environment variable")
	}

	if result, err := getEnvVar("PWD"); err != nil || result == "" {
		t.Error("expected set PWD environment variable")
	}
}
