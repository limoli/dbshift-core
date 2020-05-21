package dbshiftcore

import (
	"fmt"
	"os"
	"strconv"
)

const (
	envPathMigrations            = "DBSHIFT_ABS_FOLDER_MIGRATIONS"
	envOptionIsCreateDisabled    = "DBSHIFT_OPTION_IS_CREATE_DISABLED"
	envOptionIsDowngradeDisabled = "DBSHIFT_OPTION_IS_DOWNGRADE_DISABLED"
	envOptionIsUpgradeDisabled   = "DBSHIFT_OPTION_IS_UPGRADE_DISABLED"
)

type configuration struct {
	MigrationsPath string
	Options        configurationOptions
}

type configurationOptions struct {
	IsCreateDisabled    bool
	IsDowngradeDisabled bool
	IsUpgradeDisabled   bool
}

func getConfiguration() (*configuration, error) {
	folderMigrations, err := getEnvVar(envPathMigrations)
	if err != nil {
		return nil, err
	}

	// Check if migrations path exists
	err = checkMigrationPath(folderMigrations)
	if err != nil {
		return nil, err
	}

	// Optional values
	options, err := getOptions()
	if err != nil {
		return nil, err
	}

	return &configuration{
		MigrationsPath: folderMigrations,
		Options:        *options,
	}, nil
}

func checkMigrationPath(migrationsPath string) error {
	_, err := os.Stat(migrationsPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("migrations folder does not exist at %s", migrationsPath)
	}
	return nil
}

func getOptions() (*configurationOptions, error) {
	var err error
	options := configurationOptions{}

	if options.IsCreateDisabled, err = getBooleanOption(envOptionIsCreateDisabled); err != nil {
		return nil, err
	}

	if options.IsDowngradeDisabled, err = getBooleanOption(envOptionIsDowngradeDisabled); err != nil {
		return nil, err
	}

	if options.IsUpgradeDisabled, err = getBooleanOption(envOptionIsUpgradeDisabled); err != nil {
		return nil, err
	}

	return &options, nil
}

func getBooleanOption(envKey string) (bool, error) {
	optionEnv, err := getEnvVar(envKey)
	if err == nil {
		return strconv.ParseBool(optionEnv)
	}
	return false, nil
}

func getEnvVar(key string) (string, error) {
	var err error

	value := os.Getenv(key)
	if value == "" {
		err = fmt.Errorf("%s is not set", key)
	}

	return value, err
}
