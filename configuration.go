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
	_, err = os.Stat(folderMigrations)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("migrations folder does not exist at %s", folderMigrations)
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

func getOptions() (*configurationOptions, error) {
	options := configurationOptions{}

	optionIsCreateDisabled, err := getEnvVar(envOptionIsCreateDisabled)
	if err == nil {
		options.IsCreateDisabled, err = strconv.ParseBool(optionIsCreateDisabled)
		if err != nil {
			return nil, err
		}
	}

	optionIsDowngradeDisabled, err := getEnvVar(envOptionIsDowngradeDisabled)
	if err == nil {
		options.IsDowngradeDisabled, err = strconv.ParseBool(optionIsDowngradeDisabled)
		if err != nil {
			return nil, err
		}
	}

	optionIsUpgradeDisabled, err := getEnvVar(envOptionIsUpgradeDisabled)
	if err == nil {
		options.IsUpgradeDisabled, err = strconv.ParseBool(optionIsUpgradeDisabled)
		if err != nil {
			return nil, err
		}
	}

	return &options, nil
}

func getEnvVar(key string) (string, error) {
	var err error

	value := os.Getenv(key)
	if value == "" {
		err = fmt.Errorf("%s is not set", key)
	}

	return value, err
}
