package dbshiftcore

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var c *cmd

type dummyDbImplementation struct {
	envStatus string
}

func (db *dummyDbImplementation) GetExtension() string {
	return "txt"
}

func (db *dummyDbImplementation) GetStatus() (*Status, error) {
	s := Status{}

	statusJson := os.Getenv(db.envStatus)
	if statusJson == "" {
		return &s, nil
	}

	if err := json.Unmarshal([]byte(statusJson), &s); err != nil {
		return nil, err
	}

	return &s, nil
}

func (db *dummyDbImplementation) SetStatus(migration Migration, executionTimeInSeconds float64) error {
	s := Status{
		Type:    migration.Type,
		Version: migration.Version,
	}

	statusJson, err := json.Marshal(s)
	if err != nil {
		return err
	}

	if err := os.Setenv(db.envStatus, string(statusJson)); err != nil {
		return err
	}

	return nil
}

func (db *dummyDbImplementation) ExecuteMigration([]byte) error {
	return nil
}

func TestNewCmd_NoImplementation(t *testing.T) {
	_, err := NewCmd(nil)
	assert.NotNil(t, err, "expected missing db implementation error")
}

func TestNewCmd_NoConfiguration(t *testing.T) {
	_, err := NewCmd(new(dummyDbImplementation))
	assert.NotNil(t, err, "expected missing configuration")
}

func TestNewCmd(t *testing.T) {
	err := setConfigurationWithCustomOptions("false", "false", "false")
	assert.Nil(t, err, "expected nil error on set configuration with custom options")
	c, err = NewCmd(&dummyDbImplementation{envStatus: "DBSHIFT_DUMMY_STATUS"})
	assert.Nil(t, err, "expected nil error")
}

func TestCmd_Run_NoArgs(t *testing.T) {
	// Reset args in order to run shell in interactive mode
	os.Args = []string{}
	// Run shell in interactive mode
	c.Run()
}

func TestCmd_Run_WithArgs(t *testing.T) {
	os.Args = []string{"help"}
	c.Run()
}

func TestGetShellCommands(t *testing.T) {
	cmdList := c.getShellCommands()
	assert.Equal(t, 4, len(cmdList), "expected specific amount of commands")
}

func TestCmd_HandleStatus(t *testing.T) {
	assert.Nil(t, c.status(), "expect nil error handling status")
}

func TestCmd_HandleCreate(t *testing.T) {
	assert.Nil(t, c.create("some-migration"), "expect nil error handling create")
}

func TestCmd_HandleUpgrade(t *testing.T) {
	assert.Nil(t, c.upgrade(""), "expect nil error handling upgrade")
}

func TestCmd_HandleDowngrade(t *testing.T) {
	assert.Nil(t, c.downgrade(""), "expect nil error handling downgrade")
}

func TestCmdCreate(t *testing.T) {
	err := c.create("my-beautiful-migration")

	migrationsPath := os.Getenv("DBSHIFT_ABS_FOLDER_MIGRATIONS")
	assert.NotEmpty(t, migrationsPath, "expected valid migration path")

	t.Log(migrationsPath)
	stdout, err := exec.Command("ls", migrationsPath).Output()
	t.Log(string(stdout), err)

	assert.Nil(t, err)
}

func TestCmdStatus(t *testing.T) {
	err := c.status()
	assert.Nil(t, err)
}

func TestCmdUpgrade(t *testing.T) {
	err := c.upgrade("")
	assert.Nil(t, err, "expect nil error on upgrade")
}

func TestCmdDowngrade(t *testing.T) {
	err := c.downgrade("")
	assert.Nil(t, err, "expect nil error on downgrade")
}

// When Disabled

func TestNewCmd_Disabled(t *testing.T) {
	err := setConfigurationWithCustomOptions("true", "true", "true")
	assert.Nil(t, err, "expected nil error on set configuration with custom options")
	c, err = NewCmd(&dummyDbImplementation{envStatus: "DBSHIFT_DUMMY_STATUS"})
	assert.Nil(t, err, "expected nil error")
}

func TestCmd_Create_Disabled(t *testing.T) {
	err := c.create("")
	assert.NotNil(t, err, "expect error on create because disabled")
}

func TestCmd_Upgrade_Disabled(t *testing.T) {
	err := c.upgrade("")
	assert.NotNil(t, err, "expect error on upgrade because disabled")
}

func TestCmd_Downgrade_Disabled(t *testing.T) {
	err := c.downgrade("")
	assert.NotNil(t, err, "expect nil error on downgrade because disabled")
}

// Helpers

func setConfigurationWithCustomOptions(isCreateDisabled, isDowngradeDisabled, isUpgradeDisabled string) error {
	migrationsPath := filepath.Join("./tmp/dbshift")

	// Remove everything inside the migration path for a clean test
	if err := os.RemoveAll(migrationsPath); err != nil {
		return err
	}

	// Create temp migration path
	if err := os.MkdirAll(migrationsPath, 0777); err != nil {
		return err
	}

	err := os.Setenv("DBSHIFT_ABS_FOLDER_MIGRATIONS", migrationsPath)
	if err != nil {
		return err
	}

	if err := setIsCreateDisabled(isCreateDisabled); err != nil {
		return err
	}

	if err := setIsDowngradeDisabled(isDowngradeDisabled); err != nil {
		return err
	}

	if err := setIsUpgradeDisabled(isUpgradeDisabled); err != nil {
		return err
	}

	return nil
}

func setIsCreateDisabled(v string) error {
	return os.Setenv("DBSHIFT_OPTION_IS_CREATE_DISABLED", v)
}

func setIsDowngradeDisabled(v string) error {
	return os.Setenv("DBSHIFT_OPTION_IS_DOWNGRADE_DISABLED", v)
}

func setIsUpgradeDisabled(v string) error {
	return os.Setenv("DBSHIFT_OPTION_IS_UPGRADE_DISABLED", v)
}
