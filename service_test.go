package dbshiftcore

import (
	"testing"
	"time"
)

func TestGetMigrations(t *testing.T) {
	migrationsPath := setExistingMigrationPath(t)

	status := Status{
		Version: time.Date(2019, time.September, 1, 0, 0, 0, 0, time.Local).Format("20060102150405"),
		Type:    migrationTypeDowngrade,
	}

	migrationUpgradeList, err := getMigrations(migrationsPath, status, "", isUpgradable)
	if err != nil {
		t.Error(err)
	} else if len(migrationUpgradeList) != 2 {
		t.Errorf("unexpected counter of upgrading migrations: %d", len(migrationUpgradeList))
	}

	migrationDowngradeList, err := getMigrations(migrationsPath, status, "", isDowngradable)
	if err != nil {
		t.Error(err)
	} else if len(migrationDowngradeList) != 0 {
		t.Errorf("unexpected counter of downgrading migrations: %d ", len(migrationDowngradeList))
	}
}
