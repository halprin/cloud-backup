# cloud-backup-go
Regular back-ups to the cloud

## Install

You can get `cloud-backup` from the [releases](https://github.com/halprin/cloud-backup-go/releases)
section of this GitHub repository.  There you will find downloads for your CPU architecture.  You can then choose to
move the program to a folder contained on your [`$PATH`](https://en.wikipedia.org/wiki/PATH_(variable)).

## Run
Execute the program on your terminal.

```shell
$ cloud-backup
```

There is help built into the CLI.

You may need to allow `cloud-backup` to execute on your Mac when you first run it.  Navigate to Security & Privacy
in System Preferences to allow execution.  Also, depending on what files and folders you want to backup, you may need to
give `cloud-backup` access to certain files and folders under the Privacy tab.

## Features

### Backup

```shell
$ cloud-backup backup --help
```

Initiates a backup based on the supplied configuration file.

### Restore

```shell
$ cloud-backup restore --help
```

Restores a previously backed-up file to a specified restore location.

### List

```shell
$ cloud-backup list --help
```

List timestamps or backups at the specified timestamp.

### Install and Uninstall

```shell
$ sudo cloud-backup install --help
```

Creates a launchd daemon to execute a backup on a specified cadence.  Notice the `sudo`.  It is required because this
command creates a daemon which runs system wide and not just for a specific user.  This means `cloud-backup` runs as
`root` when executed through this manner.

```shell
$ sudo cloud-backup uninstall --help
```

Uninstalls the launchd daemon.  Notice the `sudo`.  Uninstalling a system wide daemon requires it.

## Configuration File Format

There are times that you need to pass in a path to a configuration file.  That file format is in YAML in the following
schema.

```yaml
awsCredentialConfigPath: #the full path to folder holding your shared credentials and config files; optional; if unspecified, uses the executing user's ~/.aws/ folder
aws_profile: #a profile specified in shared credentials and config files
kms_key: #a KMS key ARN; enables server-side encryption using this KMS key; optional; if unspecified, no server-side encryption is used
encryption_context: #a special value used in server-side encryption; must be a string holding JSON which represents a string-string map; optional; defaults to an empty string and is unused if the KMS key is unspecified
s3_bucket: #an S3 bucket name where the files are backed-up to
backup:
  - title: #a title to give this specific back-up file
    path: #the path to back-up
    ignore: #ignore list; optional; defaults to an empty list
      - #a string to do a simple match on which marks it for exclusion from the back-up
```
