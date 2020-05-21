package dbshiftcore

import (
	"github.com/stretchr/testify/assert"
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

	inputs := []uint{
		0,
		1,
		2,
		3,
		4,
		5,
		6,
		7,
		8,
		9,
		10,
		11,
		12,
	}

	expectedOutputs := []migrationType{
		migrationTypeDowngrade,
		migrationTypeUpgrade,
		migrationTypeDowngrade,
		migrationTypeUpgrade,
		migrationTypeDowngrade,
		migrationTypeUpgrade,
		migrationTypeDowngrade,
		migrationTypeUpgrade,
		migrationTypeDowngrade,
		migrationTypeUpgrade,
		migrationTypeDowngrade,
		migrationTypeUpgrade,
		migrationTypeDowngrade,
	}

	assert.Equal(t, len(inputs), len(expectedOutputs))

	for i := 0; i < len(inputs); i++ {
		assert.Equal(t, newMigrationTypeFromFileIndex(inputs[i]), expectedOutputs[i], "expected migration type giving file index")
	}

}

func TestNewMigrationFileName(t *testing.T) {

	inputs := []string{
		newMigrationFileName("123", "hello-world", migrationTypeUpgrade, "sql"),
		newMigrationFileName("123", "hello-world", migrationTypeDowngrade, "sql"),
		newMigrationFileName("456", "bye-world", migrationTypeUpgrade, "sql"),
		newMigrationFileName("456", "bye-world", migrationTypeDowngrade, "sql"),
		newMigrationFileName("789", "new-world", migrationTypeUpgrade, "sql"),
		newMigrationFileName("789", "new-world", migrationTypeDowngrade, "sql"),
	}

	expectedOutputs := []string{
		"123-hello-world.up.sql",
		"123-hello-world.down.sql",
		"456-bye-world.up.sql",
		"456-bye-world.down.sql",
		"789-new-world.up.sql",
		"789-new-world.down.sql",
	}

	assert.Equal(t, len(inputs), len(expectedOutputs))

	for i := 0; i < len(inputs); i++ {
		assert.Equal(t, inputs[i], expectedOutputs[i], "expected migration filename")
	}

}

func TestMigrationGetLocation(t *testing.T) {
	wd, err := os.Getwd()
	assert.Nil(t, err)

	migrationsPath := filepath.Join(wd, "example", "migrations")

	inputs := []Migration{
		newMigration("20190926154408", "hello-world", migrationTypeUpgrade, "sql"),
		newMigration("20190926154408", "hello-world", migrationTypeDowngrade, "sql"),
	}

	expectedOutputs := []string{
		filepath.Join(migrationsPath, "20190926154408-hello-world.up.sql"),
		filepath.Join(migrationsPath, "20190926154408-hello-world.down.sql"),
	}

	assert.Equal(t, len(inputs), len(expectedOutputs))

	for i := 0; i < len(inputs); i++ {
		fileName := inputs[i].getLocation(migrationsPath)
		assert.Equal(t, fileName, expectedOutputs[i], "expected same filename for migration")
	}
}

func TestMigrationFromFile(t *testing.T) {

	inputs := []string{
		"20190926154408-hello-world.down.sql",
		"20190926154408-hello-world.up.sql",
	}

	expectedOutputs := []Migration{{
		Version: "20190926154408",
		Name:    "20190926154408-hello-world.down.sql",
		Type:    migrationTypeDowngrade,
	}, {
		Version: "20190926154408",
		Name:    "20190926154408-hello-world.up.sql",
		Type:    migrationTypeUpgrade,
	}}

	assert.Equal(t, len(inputs), len(expectedOutputs))

	for i := 0; i < len(inputs); i++ {
		m, err := newMigrationFromFile(inputs[i], uint(i))
		assert.Nil(t, err)
		assert.Equal(t, m.Name, expectedOutputs[i].Name)
		assert.Equal(t, m.Version, expectedOutputs[i].Version)
		assert.Equal(t, m.Type, expectedOutputs[i].Type)
	}

}

func TestMigrationIsUpgradable(t *testing.T) {
	version := time.Now().Format("20060102150405")
	m := Migration{
		Version: version,
		Name:    "hello-world.up.sql",
		Type:    migrationTypeUpgrade,
	}

	inputs := []Status{{
		Version: time.Now().AddDate(0, 0, -1).Format("20060102150405"),
		Type:    migrationTypeDowngrade,
	}, {
		Version: time.Now().AddDate(0, 0, -2).Format("20060102150405"),
		Type:    migrationTypeDowngrade,
	}, {
		Version: time.Now().AddDate(0, 0, 1).Format("20060102150405"),
		Type:    migrationTypeDowngrade,
	}, {
		Version: m.Version,
		Type:    migrationTypeUpgrade,
	}, {
		Version: m.Version,
		Type:    migrationTypeDowngrade,
	}}

	expectedOutputs := []bool{
		true,
		true,
		false,
		false,
		true,
	}

	assert.Equal(t, len(inputs), len(expectedOutputs))

	for i := 0; i < len(inputs); i++ {
		assert.Equal(t, isUpgradable(m, inputs[i], ""), expectedOutputs[i], "expected is upgradable result")
	}
}

func TestMigrationIsDowngradable(t *testing.T) {
	version := time.Now().Format("20060102150405")
	m := Migration{
		Version: version,
		Name:    "hello-world.down.sql",
		Type:    migrationTypeDowngrade,
	}

	inputs := []Status{{
		Version: time.Now().AddDate(0, 0, -1).Format("20060102150405"),
		Type:    migrationTypeUpgrade,
	}, {
		Version: time.Now().AddDate(0, 0, -2).Format("20060102150405"),
		Type:    migrationTypeUpgrade,
	}, {
		Version: time.Now().AddDate(0, 0, 1).Format("20060102150405"),
		Type:    migrationTypeUpgrade,
	}, {
		Version: m.Version,
		Type:    migrationTypeDowngrade,
	}, {
		Version: m.Version,
		Type:    migrationTypeUpgrade,
	}}

	expectedOutputs := []bool{
		false,
		false,
		true,
		false,
		true,
	}

	assert.Equal(t, len(inputs), len(expectedOutputs))

	for i := 0; i < len(inputs); i++ {
		assert.Equal(t, isDowngradable(m, inputs[i], ""), expectedOutputs[i], "expected is downgradable result")
	}
}
