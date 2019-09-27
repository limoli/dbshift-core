package dbshiftcore

import (
	"fmt"
	"testing"
	"time"
)

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

func TestMigrationFromFile(t *testing.T) {

	const migrationName = "hello-world"
	const fileLocation = "/srv/app/migrations"

	version := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf("%s-%s", version, migrationName)

	m, err := newMigrationFromFile(fileName, 33, fileLocation)
	if err != nil {
		t.Error(err)
	}

	if m.Name != migrationName {
		t.Error(fmt.Sprintf("%s != %s", m.Name, migrationName))
	}

	if m.Version != version {
		t.Error(fmt.Sprintf("%s != %s", m.Version, version))
	}

	if m.Location != fileLocation {
		t.Error(fmt.Sprintf("%s != %s", m.Location, fileLocation))
	}

	if m.Type != migrationTypeUpgrade {
		t.Error("unexpected migration type")
	}

}

func TestMigrationIsUpgradable(t *testing.T) {

	version := time.Now().Format("20060102150405")
	m := Migration{
		Version:  version,
		Name:     "hello-world.up.sql",
		Location: "/srv/app/migrations/" + version + "-" + "hello-world.up.sql",
		Type:     migrationTypeUpgrade,
	}

	// Current version => migration
	tests := map[Status]bool{
		Status{
			Version: time.Now().AddDate(0, 0, -1).Format("20060102150405"),
			Type:    migrationTypeDowngrade,
		}: true,
		Status{
			Version: time.Now().AddDate(0, 0, -2).Format("20060102150405"),
			Type:    migrationTypeDowngrade,
		}: true,
		Status{
			Version: time.Now().AddDate(0, 0, 1).Format("20060102150405"),
			Type:    migrationTypeDowngrade,
		}: false,
		Status{
			Version: m.Version,
			Type:    migrationTypeUpgrade,
		}: false,
		Status{
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
		Version:  version,
		Name:     "hello-world.down.sql",
		Location: "/srv/app/migrations/" + version + "-" + "hello-world.down.sql",
		Type:     migrationTypeDowngrade,
	}

	// Current version => migration
	tests := map[Status]bool{
		Status{
			Version: time.Now().AddDate(0, 0, -1).Format("20060102150405"),
			Type:    migrationTypeUpgrade,
		}: false,
		Status{
			Version: time.Now().AddDate(0, 0, -2).Format("20060102150405"),
			Type:    migrationTypeUpgrade,
		}: false,
		Status{
			Version: time.Now().AddDate(0, 0, 1).Format("20060102150405"),
			Type:    migrationTypeUpgrade,
		}: true,
		Status{
			Version: m.Version,
			Type:    migrationTypeDowngrade,
		}: false,
		Status{
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
