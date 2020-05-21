package dbshiftcore

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestGetConfiguration(t *testing.T) {

	err := os.Unsetenv(envPathMigrations)
	if err != nil {
		t.Error(err)
	}

	if _, err := getConfiguration(); err == nil {
		t.Errorf("expected missing %s environment variable", envPathMigrations)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	migrationsPath := filepath.Join(wd, "example", "migrations")
	err = os.Setenv(envPathMigrations, migrationsPath)
	if err != nil {
		t.Error(err)
	}

	cfg, err := getConfiguration()
	if err != nil {
		t.Errorf("expected set %s environment variable", envPathMigrations)
	} else if cfg.MigrationsPath != migrationsPath {
		t.Errorf("expected same migration path: %s != %s", cfg.MigrationsPath, migrationsPath)
	}

}

func TestGetOptions(t *testing.T) {

	wd, err := os.Getwd()
	assert.Nil(t, err)

	migrationsPath := filepath.Join(wd, "example", "migrations")

	err = os.Setenv(envPathMigrations, migrationsPath)
	assert.Nil(t, err)

	err = os.Setenv(envOptionIsCreateDisabled, "true")
	assert.Nil(t, err)

	err = os.Setenv(envOptionIsDowngradeDisabled, "true")
	assert.Nil(t, err)

	err = os.Setenv(envOptionIsUpgradeDisabled, "true")
	assert.Nil(t, err)

	opts, err := getOptions()
	assert.Nil(t, err)
	assert.True(t, opts.IsCreateDisabled)
	assert.True(t, opts.IsDowngradeDisabled)
	assert.True(t, opts.IsUpgradeDisabled)
}

func TestGetOptions_Default(t *testing.T) {
	var err error

	err = os.Unsetenv(envPathMigrations)
	assert.Nil(t, err)

	err = os.Unsetenv(envOptionIsCreateDisabled)
	assert.Nil(t, err)

	err = os.Unsetenv(envOptionIsDowngradeDisabled)
	assert.Nil(t, err)

	err = os.Unsetenv(envOptionIsUpgradeDisabled)
	assert.Nil(t, err)

	opts, err := getOptions()
	assert.Nil(t, err)
	assert.False(t, opts.IsCreateDisabled)
	assert.False(t, opts.IsDowngradeDisabled)
	assert.False(t, opts.IsUpgradeDisabled)
}

func TestGetOptions_Error(t *testing.T) {
	wd, err := os.Getwd()
	assert.Nil(t, err)

	migrationsPath := filepath.Join(wd, "example", "migrations")

	err = os.Setenv(envPathMigrations, migrationsPath)
	assert.Nil(t, err)

	err = os.Setenv(envOptionIsCreateDisabled, "ttruuee")
	assert.Nil(t, err)

	err = os.Setenv(envOptionIsDowngradeDisabled, "ttruuee")
	assert.Nil(t, err)

	err = os.Setenv(envOptionIsUpgradeDisabled, "faalssee")
	assert.Nil(t, err)

	_, err = getOptions()
	assert.NotNil(t, err)
}

func TestGetBooleanOption(t *testing.T) {

	var b bool
	var err error

	err = os.Setenv(envOptionIsCreateDisabled, "ttruuee")
	assert.Nil(t, err)
	b, err = getBooleanOption(envOptionIsCreateDisabled)
	assert.NotNil(t, err)
	assert.False(t, b)

	err = os.Setenv(envOptionIsDowngradeDisabled, "false")
	assert.Nil(t, err)
	b, err = getBooleanOption(envOptionIsDowngradeDisabled)
	assert.Nil(t, err)
	assert.False(t, b)

	err = os.Setenv(envOptionIsUpgradeDisabled, "true")
	assert.Nil(t, err)
	b, err = getBooleanOption(envOptionIsUpgradeDisabled)
	assert.Nil(t, err)
	assert.True(t, b)
}

func TestGetEnvVar(t *testing.T) {
	if result, err := getEnvVar("unavailable_environment_variable!"); err == nil || result != "" {
		t.Error("expected missing environment variable")
	}

	if result, err := getEnvVar("PWD"); err != nil || result == "" {
		t.Error("expected set PWD environment variable")
	}
}
