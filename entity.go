package dbshiftcore

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

type iDatabase interface {
	GetExtension() string
	GetStatus() (*Status, error)
	SetStatus(migration Migration, executionTimeInSeconds float64) error
	ExecuteMigration([]byte) error
}

// Status is a structure used to identify the current (latest) migration version and type executed on database.
type Status struct {
	Version string
	Type    migrationType
}

// Migration is a structure used to group the essential information regarding the database-schema migration.
type Migration struct {
	Version string
	Name    string
	Type    migrationType
}

// Migration type
type migrationType uint

const (
	migrationTypeDowngrade migrationType = iota
	migrationTypeUpgrade
)

func (m migrationType) String() string {
	if m == migrationTypeDowngrade {
		return "down"
	}
	return "up"
}

// Migrations sort

type upgradePerspective []Migration

func (s upgradePerspective) Len() int {
	return len(s)
}
func (s upgradePerspective) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s upgradePerspective) Less(i, j int) bool {
	return s[i].Version < s[j].Version
}

type downgradePerspective []Migration

func (s downgradePerspective) Len() int {
	return len(s)
}
func (s downgradePerspective) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s downgradePerspective) Less(i, j int) bool {
	return s[i].Version > s[j].Version
}

func newMigration(version string, migrationName string, migrationType migrationType, extension string) Migration {
	return Migration{
		Version: version,
		Name:    newMigrationFileName(version, migrationName, migrationType, extension),
		Type:    migrationType,
	}
}

func newMigrationFileName(version string, migrationName string, migrationType migrationType, extension string) string {
	return fmt.Sprintf("%s-%s.%s.%s", version, migrationName, migrationType.String(), extension)
}

func (m *Migration) getLocation(migrationsPath string) string {
	return filepath.Join(migrationsPath, m.Name)
}

func newMigrationTypeFromFileIndex(fileIndex uint) migrationType {
	if fileIndex%2 == 0 {
		return migrationTypeDowngrade
	}
	return migrationTypeUpgrade
}

func newMigrationFromFile(fileName string, fileIndex uint) (*Migration, error) {
	indexDelimiter, err := getDelimiterIndexFromFileName(fileName, '-')
	if err != nil {
		return nil, err
	}

	return &Migration{
		Version: fileName[:*indexDelimiter],
		Name:    fileName,
		Type:    newMigrationTypeFromFileIndex(fileIndex),
	}, nil
}

func getDelimiterIndexFromFileName(fileName string, delimiter rune) (*int, error) {
	indexDelimiter := strings.IndexRune(fileName, '-')
	if indexDelimiter == -1 {
		return nil, errors.New("bad migration file")
	}
	return &indexDelimiter, nil
}

type migrationFilterFn func(m Migration, status Status, toInclusiveVersion string) bool

func isUpgradable(m Migration, status Status, toInclusiveVersion string) bool {

	// Only upgrading migrations
	if m.Type != migrationTypeUpgrade {
		return false
	}

	// Only migrations with version greater or equal to the current version
	if m.Version < status.Version {
		return false
	}

	// Only migrations that are not already executed
	if status.Version == m.Version && status.Type == migrationTypeUpgrade {
		return false
	}

	// If inclusive version is set, only migration with a less/equal version
	if toInclusiveVersion != "" && m.Version > toInclusiveVersion {
		return false
	}

	return true
}

func isDowngradable(m Migration, status Status, toInclusiveVersion string) bool {

	// Only downgrading migrations
	if m.Type != migrationTypeDowngrade {
		return false
	}

	// Only migrations with version greater or equal to the current version
	if m.Version > status.Version {
		return false
	}

	// Only migrations that are not already executed
	if status.Version == m.Version && status.Type == migrationTypeDowngrade {
		return false
	}

	// If inclusive version is set, only migration with a less/equal version
	if toInclusiveVersion != "" && m.Version < toInclusiveVersion {
		return false
	}

	return true
}
