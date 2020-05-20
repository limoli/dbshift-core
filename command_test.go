package dbshiftcore

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var c *cmd

type fakeDbImplementation struct{}

func (db *fakeDbImplementation) GetExtension() string {
	return "sql"
}

func (db *fakeDbImplementation) GetStatus() (*Status, error) {
	return nil, errors.New("not implemented")
}

func (db *fakeDbImplementation) SetStatus(migration Migration, executionTimeInSeconds float64) error {
	return errors.New("not implemented")
}

func (db *fakeDbImplementation) ExecuteMigration([]byte) error {
	return errors.New("not implemented")
}

func TestNewCmd_NoImplementation(t *testing.T) {
	_, err := NewCmd(nil)
	assert.NotNil(t, err, "expected missing db implementation error")
}

func TestNewCmd_NoConfiguration(t *testing.T) {
	_, err := NewCmd(new(fakeDbImplementation))
	assert.NotNil(t, err, "expected missing configuration")
}

func TestNewCmd(t *testing.T) {
	err := setConfigurationWithCustomOptions("true", "true", "false")
	assert.Nil(t, err, "expected nil error on set configuration with custom options")

	c, err = NewCmd(new(fakeDbImplementation))
	assert.Nil(t, err, "expected nil error")
}

func TestCmd_Run(t *testing.T) {
	// Reset args in order to run shell in interactive mode
	os.Args = []string{}
	// Run shell in interactive mode
	c.Run()
}

func TestGetShellCommands(t *testing.T) {
	cmdList := c.getShellCommands()
	assert.Equal(t, 4, len(cmdList), "expected specific amount of commands")
}

func TestCmdCreate(t *testing.T) {
	err := c.create("my-beautiful-migration")
	assert.Error(t, err, "expect error since no db implementation has been passed - db extension is used to generate the migration filename")
}

func TestCmdStatus(t *testing.T) {
	err := c.status()
	assert.Error(t, err, "expect error since no db implementation has been passed")
}

func TestCmdUpgrade(t *testing.T) {
	err := c.upgrade("")
	assert.Error(t, err, "expect error since upgrading is disabled")

}

func TestCmdDowngrade(t *testing.T) {
	err := c.downgrade("")
	assert.Error(t, err, "expect error since downgrading is disabled")
}

func setConfigurationWithCustomOptions(isCreateDisabled, isDowngradeDisabled, isUpgradeDisabled string) error {
	migrationsPath := os.TempDir()

	err := os.Setenv("DBSHIFT_ABS_FOLDER_MIGRATIONS", migrationsPath)
	if err != nil {
		return err
	}

	err = os.Setenv("DBSHIFT_OPTION_IS_CREATE_DISABLED", isCreateDisabled)
	if err != nil {
		return err
	}

	err = os.Setenv("DBSHIFT_OPTION_IS_DOWNGRADE_DISABLED", isDowngradeDisabled)
	if err != nil {
		return err
	}

	err = os.Setenv("DBSHIFT_OPTION_IS_UPGRADE_DISABLED", isUpgradeDisabled)
	if err != nil {
		return err
	}

	return nil
}
