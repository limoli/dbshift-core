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
		Version: time.Now().AddDate(0, 0, -1).Format("20060102150405"),
		Type:    migrationTypeDowngrade,
	}

	migrationUpgradeList, err := getMigrations(migrationsFolder, status, "", isUpgradable)
	if err != nil {
		t.Error(err)
	} else if len(migrationUpgradeList) != 2 {
		t.Error("unexpected counter of upgrading migrations")
	}

	migrationDowngradeList, err := getMigrations(migrationsFolder, status, "", isDowngradable)
	if err != nil {
		t.Error(err)
	} else if len(migrationDowngradeList) != 0 {
		t.Error("unexpected counter of downgrading migrations")
	}

}
