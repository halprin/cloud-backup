package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type BackupConfiguration struct {
	AwsProfile        string       `yaml:"aws_profile"`
	KmsKey            string       `yaml:"kms_key"`
	EncryptionContext string       `yaml:"encryption_context"`
	S3Bucket          string       `yaml:"s3_bucket"`
	IntermediatePath  string       `yaml:"intermediate_path"`
	Backup            []BackupFile
}

type BackupFile struct {
	Title  string
	Path   string
	Ignore []string
}

var backupConfig *BackupConfiguration

func BackupConfig() (BackupConfiguration, error) {
	if backupConfig == nil {
		tempBackupConfig, err := parse(os.Args[1])
		if err != nil {
			return BackupConfiguration{}, err
		}

		backupConfig = &tempBackupConfig
	}

	return *backupConfig, nil
}

func parse(filePath string) (BackupConfiguration, error) {
	log.Println("Parsing backup config")

	configRawData, err := os.ReadFile(filePath)
	if err != nil {
		return BackupConfiguration{}, err
	}

	internalBackupConfig := BackupConfiguration{}
	err = yaml.Unmarshal(configRawData, &internalBackupConfig)
	return internalBackupConfig, err
}
