package actions

import (
	"github.com/halprin/cloud-backup-go/config"
	"github.com/halprin/cloud-backup-go/crypt"
	"github.com/halprin/cloud-backup-go/transfer"
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

	downloader, err := transfer.NewDownloader(overallConfig, timestamp, backupFile)
	if err != nil {
		return err
	}

	downloadReader, err := downloader.Download()
	if err != nil {
		return err
	}

	outputFilePath := filepath.Join(restorePath, backupFile + ".tar.gz")
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	decryptor := crypt.NewDecryptor(downloadReader, outputFile, overallConfig)

	err = decryptor.Decrypt()
	if err != nil {
		return err
	}

	log.Println("Restoring file complete")
	return nil
}
