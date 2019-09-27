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

func TestMigrationGetFilename(t *testing.T) {
	tests := map[Migration]string{
		{
			Version: "123",
			Name:    "hello-world.up.sql",
			Type:    migrationTypeUpgrade,
		}: "123-hello-world.up.sql",
		{
			Version: "123",
			Name:    "hello-world.down.sql",
			Type:    migrationTypeDowngrade,
		}: "123-hello-world.down.sql",
	}

	for k, v := range tests {
		fileName := k.getFileName()
		if fileName != v {
			t.Errorf("unexpected migration filename %s instead of %s", fileName, v)
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
		{
			Version: "123",
			Name:    "hello-world.up.sql",
			Type:    migrationTypeUpgrade,
		}: filepath.Join(migrationsPath, "123-hello-world.up.sql"),
		{
			Version: "123",
			Name:    "hello-world.down.sql",
			Type:    migrationTypeDowngrade,
		}: filepath.Join(migrationsPath, "123-hello-world.down.sql"),
	}

	for k, v := range tests {
		fileName := k.getLocation(migrationsPath)
		if fileName != v {
			t.Errorf("unexpected migration filename %s instead of %s", fileName, v)
		}
	}
}

func TestMigrationFromFile(t *testing.T) {
	const migrationName = "hello-world"

	version := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf("%s-%s", version, migrationName)

	m, err := newMigrationFromFile(fileName, 33)
	if err != nil {
		t.Error(err)
	}

	if m.Name != migrationName {
		t.Error(fmt.Sprintf("%s != %s", m.Name, migrationName))
	}

	if m.Version != version {
		t.Error(fmt.Sprintf("%s != %s", m.Version, version))
	}

	if m.Type != migrationTypeUpgrade {
		t.Error("unexpected migration type")
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
