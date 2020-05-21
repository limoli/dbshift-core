package dbshiftcore

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestGetConfiguration_MissingMigrationPath(t *testing.T) {
	err := os.Unsetenv(envPathMigrations)
	assert.Nil(t, err)

	migrationsPath := os.Getenv(envPathMigrations)
	assert.Nil(t, err)

	errCheck := checkMigrationPath(migrationsPath)
	assert.NotNil(t, errCheck)

	_, errOptions := getOptions()
	assert.Nil(t, errOptions)

	cfg, err := getConfiguration()
	assert.NotNil(t, err)
	assert.Nil(t, cfg)
}

func TestGetConfiguration_UnexistingMigrationPath(t *testing.T) {
	migrationsPath := setUnexistingMigrationPath(t)

	errCheck := checkMigrationPath(migrationsPath)
	assert.NotNil(t, errCheck)

	_, errOptions := getOptions()
	assert.Nil(t, errOptions)

	cfg, err := getConfiguration()
	assert.NotNil(t, err)
	assert.Nil(t, cfg)
}

func TestGetConfiguration_InvalidOptions(t *testing.T) {
	migrationsPath := setExistingMigrationPath(t)

	errCheck := checkMigrationPath(migrationsPath)
	assert.Nil(t, errCheck)

	setInvalidOption(t)
	_, errOptions := getOptions()
	assert.NotNil(t, errOptions)

	cfg, err := getConfiguration()
	assert.NotNil(t, err)
	assert.Nil(t, cfg)
}

func TestGetConfiguration_Present(t *testing.T) {
	migrationsPath := setExistingMigrationPath(t)

	errCheck := checkMigrationPath(migrationsPath)
	assert.Nil(t, errCheck)

	unsetInvalidOption(t)
	_, errOptions := getOptions()
	assert.Nil(t, errOptions)

	cfg, err := getConfiguration()
	assert.Nil(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, cfg.MigrationsPath, migrationsPath)
}

func TestCheckMigrationPath_Missing(t *testing.T) {
	migrationsPath := setUnexistingMigrationPath(t)
	err := checkMigrationPath(migrationsPath)
	assert.NotNil(t, err, "expected unavailable migration path")
}

func TestCheckMigrationPath_Present(t *testing.T) {
	migrationsPath := setExistingMigrationPath(t)
	err := checkMigrationPath(migrationsPath)
	assert.Nil(t, err, "expected available migration path")
}

func TestGetOptions(t *testing.T) {
	var err error
	_ = setExistingMigrationPath(t)

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

func TestGetOptions_Cases(t *testing.T) {
	setExistingMigrationPath(t)

	type test struct {
		inputValue           string
		isOptionsError       bool
		expectedBooleanValue bool
	}

	tests := []test{{
		inputValue:           "not a boolean, sure! ;)",
		isOptionsError:       true,
		expectedBooleanValue: false,
	}, {
		inputValue:           "true",
		isOptionsError:       false,
		expectedBooleanValue: true,
	}, {
		inputValue:           "false",
		isOptionsError:       false,
		expectedBooleanValue: false,
	}}

	for _, v := range tests {
		err := os.Setenv(envOptionIsCreateDisabled, v.inputValue)
		assert.Nil(t, err)
		opts, err := getOptions()
		assert.Equal(t, v.isOptionsError, err != nil)
		assert.Equal(t, v.expectedBooleanValue, opts != nil && opts.IsCreateDisabled)
	}

	for _, v := range tests {
		err := os.Setenv(envOptionIsDowngradeDisabled, v.inputValue)
		assert.Nil(t, err)
		opts, err := getOptions()
		assert.Equal(t, v.isOptionsError, err != nil)
		assert.Equal(t, v.expectedBooleanValue, opts != nil && opts.IsDowngradeDisabled)
	}

	for _, v := range tests {
		err := os.Setenv(envOptionIsUpgradeDisabled, v.inputValue)
		assert.Nil(t, err)
		opts, err := getOptions()
		assert.Equal(t, v.isOptionsError, err != nil)
		assert.Equal(t, v.expectedBooleanValue, opts != nil && opts.IsUpgradeDisabled)
	}
}

func TestGetBooleanOption_IsCreateDisabled(t *testing.T) {
	tests := map[string]expectedOutput{
		"ttruuee": {b: false, hasError: true},
		"false":   {b: false, hasError: false},
		"true":    {b: true, hasError: false},
	}
	testGetBooleanOption(t, envOptionIsCreateDisabled, tests)
}

func TestGetBooleanOption_IsDowngradeDisabled(t *testing.T) {
	tests := map[string]expectedOutput{
		"ttruuee": {b: false, hasError: true},
		"false":   {b: false, hasError: false},
		"true":    {b: true, hasError: false},
	}
	testGetBooleanOption(t, envOptionIsDowngradeDisabled, tests)
}

func TestGetBooleanOption_IsUpgradeDisabled(t *testing.T) {
	tests := map[string]expectedOutput{
		"ttruuee": {b: false, hasError: true},
		"false":   {b: false, hasError: false},
		"true":    {b: true, hasError: false},
	}
	testGetBooleanOption(t, envOptionIsUpgradeDisabled, tests)
}

func TestGetEnvVar(t *testing.T) {
	if result, err := getEnvVar("unavailable_environment_variable!"); err == nil || result != "" {
		t.Error("expected missing environment variable")
	}

	if result, err := getEnvVar("PWD"); err != nil || result == "" {
		t.Error("expected set PWD environment variable")
	}
}

// Helpers

func setExistingMigrationPath(t *testing.T) string {
	wd, err := os.Getwd()
	assert.Nil(t, err)
	migrationsPath := filepath.Join(wd, "example", "migrations")
	err = os.Setenv(envPathMigrations, migrationsPath)
	assert.Nil(t, err)
	return migrationsPath
}

func setUnexistingMigrationPath(t *testing.T) string {
	wd, err := os.Getwd()
	assert.Nil(t, err)
	migrationsPath := filepath.Join(wd, "example", "unexisting", "folder")
	err = os.Setenv(envPathMigrations, migrationsPath)
	assert.Nil(t, err)
	return migrationsPath
}

func setInvalidOption(t *testing.T) {
	err := os.Setenv(envOptionIsUpgradeDisabled, "true!")
	assert.Nil(t, err)
}

func unsetInvalidOption(t *testing.T) {
	err := os.Unsetenv(envOptionIsUpgradeDisabled)
	assert.Nil(t, err)
}

type expectedOutput struct {
	b        bool
	hasError bool
}

func testGetBooleanOption(t *testing.T, envKey string, tests map[string]expectedOutput) {
	for envValue, expectedOutput := range tests {
		err := os.Setenv(envKey, envValue)
		assert.Nil(t, err)
		b, err := getBooleanOption(envKey)
		assert.Equal(t, expectedOutput.b, b, "expected same boolean value")
		assert.Equal(t, expectedOutput.hasError, err != nil, "expected same error value")
	}
}
