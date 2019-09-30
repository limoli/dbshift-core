package dbshiftcore

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestMigrationType_String(t *testing.T) {
	if migrationTypeUpgrade.String() != "up" {
		t.Error("expected up as upgrade type to string")
	}

	if migrationTypeDowngrade.String() != "down" {
		t.Error("expected up as downgrade type to string")
	}
}

func TestMigrationTypeSort(t *testing.T) {

	migrationList := []Migration{
		{
			Version: "123",
			Name:    "hello-world.up.sql",
			Type:    migrationTypeUpgrade,
		},
		{
			Version: "456",
			Name:    "bye-world.up.sql",
			Type:    migrationTypeUpgrade,
		},
	}

	upgradePerspective := upgradePerspective(migrationList)
	if upgradePerspective.Len() != 2 {
		t.Error("unexpected length of upgrading migrations")
	}

	if upgradePerspective.Less(1, 0) {
		t.Error("unexpected less implementation for upgrading migrations")
	}

	downgradePerspective := downgradePerspective(migrationList)
	if downgradePerspective.Len() != 2 {
		t.Error("unexpected length of downgrading migrations")
	}

	if downgradePerspective.Less(0, 1) {
		t.Error("unexpected less implementation for downgrading migrations")
	}

}

func TestNewMigration(t *testing.T) {
	tests := map[uint]migrationType{
		0:  migrationTypeDowngrade,
		1:  migrationTypeUpgrade,
		2:  migrationTypeDowngrade,
		3:  migrationTypeUpgrade,
		4:  migrationTypeDowngrade,
		5:  migrationTypeUpgrade,
		6:  migrationTypeDowngrade,
		7:  migrationTypeUpgrade,
		8:  migrationTypeDowngrade,
		9:  migrationTypeUpgrade,
		10: migrationTypeDowngrade,
	}

	for k, v := range tests {
		if newMigrationTypeFromFileIndex(k) != v {
			t.Error("unexpected migration type giving file index")
		}
	}
}

func TestNewMigrationFileName(t *testing.T) {

	tests := map[string]string{
		newMigrationFileName("123", "hello-world", migrationTypeUpgrade, "sql"):   "123-hello-world.up.sql",
		newMigrationFileName("123", "hello-world", migrationTypeDowngrade, "sql"): "123-hello-world.down.sql",
		newMigrationFileName("456", "bye-world", migrationTypeUpgrade, "sql"):     "456-bye-world.up.sql",
		newMigrationFileName("456", "bye-world", migrationTypeDowngrade, "sql"):   "456-bye-world.down.sql",
		newMigrationFileName("789", "new-world", migrationTypeUpgrade, "sql"):     "789-new-world.up.sql",
		newMigrationFileName("789", "new-world", migrationTypeDowngrade, "sql"):   "789-new-world.down.sql",
	}

	for k, v := range tests {
		if k != v {
			t.Errorf("unexpected migration filename %s instead of %s", k, v)
		}
	}
}

func TestMigrationGetLocation(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	migrationsPath := filepath.Join(wd, "example", "migrations")

	tests := map[Migration]string{
		newMigration("20190926154408", "hello-world", migrationTypeUpgrade, "sql"):
		filepath.Join(migrationsPath, "20190926154408-hello-world.up.sql"),
		newMigration("20190926154408", "hello-world", migrationTypeDowngrade, "sql"):
		filepath.Join(migrationsPath, "20190926154408-hello-world.down.sql"),
	}

	for k, v := range tests {
		fileName := k.getLocation(migrationsPath)
		if fileName != v {
			t.Errorf("unexpected migration filename %s instead of %s", fileName, v)
		}
	}
}

func TestMigrationFromFile(t *testing.T) {

	tests := map[string]Migration{
		"20190926154408-hello-world.down.sql": {
			Version: "20190926154408",
			Name:    "20190926154408-hello-world.down.sql",
			Type:    migrationTypeDowngrade,
		},
		"20190926154408-hello-world.up.sql": {
			Version: "20190926154408",
			Name:    "20190926154408-hello-world.up.sql",
			Type:    migrationTypeUpgrade,
		},
	}

	var fileIndex uint = 0
	for k, v := range tests {
		m, err := newMigrationFromFile(k, fileIndex)
		if err != nil {
			t.Error(err)
		} else if m.Name != v.Name {
			t.Error(fmt.Errorf("unexpected name %v instead of %v", m.Name, v.Name))
		} else if m.Version != v.Version {
			t.Error(fmt.Errorf("unexpected version %v instead of %v", m.Version, v.Version))
		} else if m.Type != v.Type {
			t.Error(fmt.Errorf("unexpected type %v instead of %v", m.Type, v.Type))
		}
		fileIndex++
	}

}

func TestMigrationIsUpgradable(t *testing.T) {
	version := time.Now().Format("20060102150405")
	m := Migration{
		Version: version,
		Name:    "hello-world.up.sql",
		Type:    migrationTypeUpgrade,
	}

	// Current version => migration
	tests := map[Status]bool{
		{
			Version: time.Now().AddDate(0, 0, -1).Format("20060102150405"),
			Type:    migrationTypeDowngrade,
		}: true,
		{
			Version: time.Now().AddDate(0, 0, -2).Format("20060102150405"),
			Type:    migrationTypeDowngrade,
		}: true,
		{
			Version: time.Now().AddDate(0, 0, 1).Format("20060102150405"),
			Type:    migrationTypeDowngrade,
		}: false,
		{
			Version: m.Version,
			Type:    migrationTypeUpgrade,
		}: false,
		{
			Version: m.Version,
			Type:    migrationTypeDowngrade,
		}: true,
	}

	for k, v := range tests {
		if isUpgradable(m, k, "") != v {
			t.Error("unexpected is upgradable result")
		}
	}
}

func TestMigrationIsDowngradable(t *testing.T) {
	version := time.Now().Format("20060102150405")
	m := Migration{
		Version: version,
		Name:    "hello-world.down.sql",
		Type:    migrationTypeDowngrade,
	}

	// Current version => migration
	tests := map[Status]bool{
		{
			Version: time.Now().AddDate(0, 0, -1).Format("20060102150405"),
			Type:    migrationTypeUpgrade,
		}: false,
		{
			Version: time.Now().AddDate(0, 0, -2).Format("20060102150405"),
			Type:    migrationTypeUpgrade,
		}: false,
		{
			Version: time.Now().AddDate(0, 0, 1).Format("20060102150405"),
			Type:    migrationTypeUpgrade,
		}: true,
		{
			Version: m.Version,
			Type:    migrationTypeDowngrade,
		}: false,
		{
			Version: m.Version,
			Type:    migrationTypeUpgrade,
		}: true,
	}

	for k, v := range tests {
		if isDowngradable(m, k, "") != v {
			t.Error("unexpected is upgradable result")
		}
	}
}
