package restore

import (
	"github.com/halprin/cloud-backup-go/config"
	"github.com/halprin/cloud-backup-go/crypt"
	"log"
	"os"
	"path/filepath"
)

func Restore(configFilePath string, timestamp string, backupFile string, restorePath string) error {
	overallConfig, err := config.New(configFilePath)
	if err != nil {
		return err
	}

	log.Printf("Restoring file %s from %s to %s", backupFile, timestamp, restorePath)

	sourceReader, err := getSourceReader(overallConfig, timestamp, backupFile)
	if err != nil {
		return err
	}

	outputFilePath := filepath.Join(restorePath, backupFile + ".tar.gz")
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	decryptor := crypt.NewDecryptor(sourceReader, outputFile, overallConfig)

	err = decryptor.Decrypt()
	if err != nil {
		return err
	}

	log.Println("Restoring file complete")
	return nil
}
