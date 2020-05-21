package dbshiftcore

import (
	"errors"
	"fmt"
	"github.com/abiosoft/ishell"
	"io/ioutil"
	"os"
	"sort"
	"time"
)

type cmd struct {
	cfg configuration
	db  iDatabase
}

// NewCmd create a shell-commander object based on database interface and environmental configuration.
func NewCmd(db iDatabase) (*cmd, error) {

	// Check db implementation
	if db == nil {
		return nil, errors.New("missing db implementation")
	}

	// Get configuration via environment
	cfg, err := getConfiguration()
	if err != nil {
		return nil, fmt.Errorf("bad configuration: %s", err)
	}

	return &cmd{cfg: *cfg, db: db}, nil
}

// Run is used to execute the shell-commander.
func (c *cmd) Run() {

	// Run shell
	shell := ishell.New()

	commands := c.getShellCommands()
	for k := range commands {
		shell.AddCmd(commands[k])
	}

	if len(os.Args) > 1 {
		if err := shell.Process(os.Args[1:]...); err != nil {
			PrintFailure(err.Error())
			os.Exit(1)
		}
	} else {
		shell.Run()
	}
}

func (c *cmd) getShellCommands() []*ishell.Cmd {
	return []*ishell.Cmd{{
		Name:     "status",
		LongHelp: "It returns the current status of database along migrations.",
		Func:     c.handleStatus,
	}, {
		Name:     "create",
		Help:     "create <entity-name>",
		LongHelp: "It creates a entity with name.",
		Func:     c.handleCreate,
	}, {
		Name:     "upgrade",
		Help:     "upgrade [toInclusiveVersion]",
		LongHelp: "It upgrades all the migrations. If toInclusiveId is set, it upgrades all the migrations till that version.",
		Func:     c.handleUpgrade,
	}, {
		Name:     "downgrade",
		Help:     "downgrade [toInclusiveVersion]",
		LongHelp: "It downgrades all the migrations. If toInclusiveId is set, it downgrades all the migrations till that version.",
		Func:     c.handleDowngrade,
	}}
}

func (c *cmd) handleStatus(ctx *ishell.Context) {
	if err := c.status(); err != nil {
		PrintFailure(err.Error())
	}
}

func (c *cmd) handleCreate(ctx *ishell.Context) {
	if len(ctx.Args) != 1 {
		PrintFailure("missing entity name")
		return
	}
	name := ctx.Args[0]
	if err := c.create(name); err != nil {
		PrintFailure(err.Error())
	}
}

func (c *cmd) handleUpgrade(ctx *ishell.Context) {
	var endVersion string
	if len(ctx.Args) == 1 {
		endVersion = ctx.Args[0]
	}
	if err := c.upgrade(endVersion); err != nil {
		PrintFailure(err.Error())
	}
}

func (c *cmd) handleDowngrade(ctx *ishell.Context) {
	var endVersion string
	if len(ctx.Args) == 1 {
		endVersion = ctx.Args[0]
	}
	if err := c.downgrade(endVersion); err != nil {
		PrintFailure(err.Error())
	}
}

func (c *cmd) create(migrationName string) error {
	// Check option
	if c.cfg.Options.IsCreateDisabled {
		return errors.New("migration creating is disabled from options")
	}

	// Ensure both downgrading and upgrading migrations share the same version
	version := time.Now().Format("20060102150405")
	dbExt := c.db.GetExtension()

	// Write downgrade file
	migrationDowngrade := newMigration(version, migrationName, migrationTypeDowngrade, dbExt)
	if err := ioutil.WriteFile(migrationDowngrade.getLocation(c.cfg.MigrationsPath), nil, 0664); err != nil {
		return err
	}

	// Write upgrade file
	migrationUpgrade := newMigration(version, migrationName, migrationTypeUpgrade, dbExt)
	if err := ioutil.WriteFile(migrationUpgrade.getLocation(c.cfg.MigrationsPath), nil, 0664); err != nil {
		return err
	}

	return nil
}

func (c *cmd) upgrade(toInclusiveVersion string) error {
	// Check option
	if c.cfg.Options.IsUpgradeDisabled {
		return errors.New("migration upgrading is disabled from options")
	}

	// Get current version
	status, err := c.db.GetStatus()
	if err != nil {
		return err
	}

	// Get migrations eligible to upgrade
	migrationList, err := getMigrations(c.cfg.MigrationsPath, *status, toInclusiveVersion, isUpgradable)
	if err != nil {
		return err
	}

	// Sort for execution
	sort.Sort(upgradePerspective(migrationList))

	// Execute migrations
	return c.execMigrations(migrationList)
}

func (c *cmd) downgrade(toInclusiveVersion string) error {
	// Check option
	if c.cfg.Options.IsDowngradeDisabled {
		return errors.New("migration downgrading is disabled from options")
	}

	// Get current version
	status, err := c.db.GetStatus()
	if err != nil {
		return err
	}

	// Get migrations eligible to downgrade
	migrationList, err := getMigrations(c.cfg.MigrationsPath, *status, toInclusiveVersion, isDowngradable)
	if err != nil {
		return err
	}

	// Sort for execution
	sort.Sort(downgradePerspective(migrationList))

	// Execute migrations
	return c.execMigrations(migrationList)
}

func (c *cmd) execMigrations(migrationList []Migration) error {
	for _, m := range migrationList {

		// Read migration file
		data, err := ioutil.ReadFile(m.getLocation(c.cfg.MigrationsPath))
		if err != nil {
			return err
		}

		// Execute migration
		timeStart := time.Now()
		if err := c.db.ExecuteMigration(data); err != nil {
			return err
		}

		execTimeInSeconds := time.Since(timeStart).Seconds()
		if err := c.db.SetStatus(m, execTimeInSeconds); err != nil {
			return err
		}

		PrintSuccess("Migration %s has been executed in %v seconds", m.Name, execTimeInSeconds)
	}

	return nil
}

func (c *cmd) status() error {

	// Get current version
	status, err := c.db.GetStatus()
	if err != nil {
		return err
	}

	// Get migrations eligible to upgrade
	migrationUpgradeList, err := getMigrations(c.cfg.MigrationsPath, *status, "", isUpgradable)
	if err != nil {
		return err
	}

	sort.Sort(upgradePerspective(migrationUpgradeList))

	// Get migrations eligible to downgrade
	migrationDowngradeList, err := getMigrations(c.cfg.MigrationsPath, *status, "", isDowngradable)
	if err != nil {
		return err
	}

	sort.Sort(downgradePerspective(migrationDowngradeList))

	fmt.Println("Migrations to upgrade")
	for _, m := range migrationUpgradeList {
		fmt.Println(m.Name)
	}

	fmt.Println("Migrations to downgrade")
	for _, m := range migrationDowngradeList {
		fmt.Println(m.Name)
	}

	return nil
}
