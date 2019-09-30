package dbshiftcore

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGetMigrations(t *testing.T) {

	wd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	migrationsFolder := filepath.Join(wd, "example", "migrations")
	status := Status{
		Version: time.Date(2019, time.September, 1, 0, 0, 0, 0, time.Local).Format("20060102150405"),
		Type:    migrationTypeDowngrade,
	}

	migrationUpgradeList, err := getMigrations(migrationsFolder, status, "", isUpgradable)
	if err != nil {
		t.Error(err)
	} else if len(migrationUpgradeList) != 2 {
		t.Errorf("unexpected counter of upgrading migrations: %d", len(migrationUpgradeList))
	}

	migrationDowngradeList, err := getMigrations(migrationsFolder, status, "", isDowngradable)
	if err != nil {
		t.Error(err)
	} else if len(migrationDowngradeList) != 0 {
		t.Errorf("unexpected counter of downgrading migrations: %d ", len(migrationDowngradeList))
	}

}
