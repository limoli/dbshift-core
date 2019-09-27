package dbshiftcore

import (
	"errors"
	"os"
	"testing"
)

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

func TestNewCmd(t *testing.T) {
	if _, err := NewCmd(nil); err == nil {
		t.Error("expected missing db implementation error")
	}

	if _, err := NewCmd(new(fakeDbImplementation)); err == nil {
		t.Error("expected missing configuration")
	}
}

func TestGetShellCommands(t *testing.T) {
	err := setConfigurationWithCustomOptions("false", "false", "false")
	if err != nil {
		t.Error(err)
	}

	cmd, err := NewCmd(new(fakeDbImplementation))
	if err != nil {
		t.Error(err)
	}

	cmds := cmd.getShellCommands()
	if len(cmds) != 4 {
		t.Errorf("unexpected amount of commands %d instead of %d", len(cmds), 4)
	}
}

func TestCmdCreate(t *testing.T) {
	err := setConfigurationWithCustomOptions("true", "true", "false")
	if err != nil {
		t.Error(err)
	}

	cmd, err := NewCmd(new(fakeDbImplementation))
	if err != nil {
		t.Error(err)
	}

	if err := cmd.create("my-beautiful-migration"); err == nil {
		t.Error("expect error since no db implementation has been passed - db extension is used to generate the migration filename")
	}
}

func TestCmdStatus(t *testing.T) {
	err := setConfigurationWithCustomOptions("false", "false", "false")
	if err != nil {
		t.Error(err)
	}

	cmd, err := NewCmd(new(fakeDbImplementation))
	if err != nil {
		t.Error(err)
	}

	if err := cmd.status(); err == nil {
		t.Error("expect error since no db implementation has been passed")
	}
}

func TestCmdUpgrade(t *testing.T) {
	err := setConfigurationWithCustomOptions("false", "false", "true")
	if err != nil {
		t.Error(err)
	}

	cmd, err := NewCmd(new(fakeDbImplementation))
	if err != nil {
		t.Error(err)
	}

	if err := cmd.upgrade(""); err == nil {
		t.Error("expect error since upgrading is disabled")
	}
}

func TestCmdDowngrade(t *testing.T) {
	err := setConfigurationWithCustomOptions("false", "true", "false")
	if err != nil {
		t.Error(err)
	}

	cmd, err := NewCmd(new(fakeDbImplementation))
	if err != nil {
		t.Error(err)
	}

	if err := cmd.downgrade(""); err == nil {
		t.Error("expect error since downgrading is disabled")
	}
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
