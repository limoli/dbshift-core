package dbshiftcore

import (
	"os"
	"path/filepath"
)

func getMigrations(migrationsPath string, status Status, toInclusiveVersion string, filterFn migrationFilterFn) ([]Migration, error) {
	var migrationList []Migration
	var fileIndex uint

	err := filepath.Walk(migrationsPath, func(path string, info os.FileInfo, err error) error {

		if info == nil {
			return nil
		}

		fileName := info.Name()

		// Exclude directories and hidden files
		if info.IsDir() || fileName[0] == '.' {
			return nil
		}

		migrationObj, err := newMigrationFromFile(fileName, fileIndex)
		if err != nil {
			return err
		}

		if filterFn(*migrationObj, status, toInclusiveVersion) {
			migrationList = append(migrationList, *migrationObj)
		}

		fileIndex++

		return nil
	})

	return migrationList, err
}
