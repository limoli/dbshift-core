# DbShift Core

DbShift Core provides simple and light logic for the management of **database-schema migrations**.
You will be able to create migrations, check the current db status, decide to upgrade or downgrade easily.
It can be easily implemented with specific database clients.

[![GoDoc](https://godoc.org/limoli/dbshift-core?status.svg)](https://godoc.org/github.com/limoli/dbshift-core)
[![Build Status](https://travis-ci.org/limoli/dbshift-core.svg?branch=master)](https://travis-ci.org/limoli/dbshift-core)
[![Go Report Card](https://goreportcard.com/badge/github.com/limoli/dbshift-core)](https://goreportcard.com/report/github.com/limoli/dbshift-core)
[![Maintainability](https://api.codeclimate.com/v1/badges/0b1f0599ef4c4a763953/maintainability)](https://codeclimate.com/github/limoli/dbshift-core/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/0b1f0599ef4c4a763953/test_coverage)](https://codeclimate.com/github/limoli/dbshift-core/test_coverage)
[![License](http://img.shields.io/badge/license-mit-blue.svg)](https://raw.githubusercontent.com/github.com/limoli/dbshift-core/LICENSE)

## Install

DbShift Core aims to provide logic and **not an installable command**.
You can use a **DbShift Client** implementation according to your database:
- [MySQL client](https://github.com/limoli/dbshift-cli-mysql)

## Commands

Set your [configuration](#configuration)

#### Create migration 
It creates two files (`$timestamp.down.sql` and `$timestamp.up.sql`) at your migrations folder.
```bash
dbshift create my-migration-description
```

#### Status   
Check status of your migrations.
```bash
dbshift status
```
#### Upgrade
Upgrade migrations.
```bash
dbshift upgrade
```
```bash
dbshift upgrade <toInclusiveMigrationVersion>
```

#### Downgrade
Downgrade migrations.    
```bash
dbshift downgrade
```
```bash
dbshift downgrade <toInclusiveMigrationVersion>
```

## Configuration

| Key                                   | Description                                        | Value example              |
|---                                    |---                                                 |---                         |
|`DBSHIFT_ABS_FOLDER_MIGRATIONS`        | Where migrations are created and stored.           | `/srv/app/migrations`      |
|`DBSHIFT_OPTION_IS_CREATE_DISABLED`    | Disable create command (useful on production).     | `true` / `false` (default) |
|`DBSHIFT_OPTION_IS_DOWNGRADE_DISABLED` | Disable downgrade command (useful on production).  | `true` / `false` (default) |
|`DBSHIFT_OPTION_IS_UPGRADE_DISABLED`   | Disable upgrade command (useful on production).    | `true` / `false` (default) |	

This configuration represents the basic configuration for the DbShift Core.
More configurations can be offered by the single DbShift Client.

## Write good migrations

1. Queries must be database name **agnostic**
2. [SRP](https://en.wikipedia.org/wiki/Single_responsibility_principle) according to your description
3. Write both upgrade and downgrade migrations 

## Exit codes

The following error-codes interval is reserved for core usage: `[1, 90]`.

| Code      | Description                                                           |
| ---       | ---                                                                   |
| `1`       | When no command is passed in the no-interactive mode.                 |

## Client implementation

#### Exit codes

The client implementation interval is `[100,255]`.

#### Environment

The environment variables must have the following prefix: `DBSHIFT_CLI_<DBTYPE>`.

Example for `MySQL`: 
- Prefix: `DBSHIFT_CLI_MYSQL`
- Variables: `DBSHIFT_CLI_MYSQL_X`, `DBSHIFT_CLI_MYSQL_Y`, `DBSHIFT_CLI_MYSQL_Z`
