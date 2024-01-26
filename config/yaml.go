package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type BackupConfiguration struct {
	AwsCredentialConfigPath string                    `yaml:"awsCredentialConfigPath"`
	AwsProfile              string                    `yaml:"aws_profile"`
	KmsKey                  string                    `yaml:"kms_key"`
	EncryptionContext       string                    `yaml:"encryption_context"`
	S3Bucket                string                    `yaml:"s3_bucket"`
	BackupFiles             []BackupFileConfiguration `yaml:"backup"`
}

type BackupFileConfiguration struct {
	Title  string
	Path   string
	Ignore []string
}

var backupConfig *BackupConfiguration

func BackupConfig() (BackupConfiguration, error) {
	if backupConfig == nil {
		return BackupConfiguration{}, fmt.Errorf("backup config not initialized first")
	}

	return *backupConfig, nil
}

func New(configFilePath string) (BackupConfiguration, error) {
	tempBackupConfig, err := parse(configFilePath)
	if err != nil {
		return BackupConfiguration{}, err
	}

	backupConfig = &tempBackupConfig

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
