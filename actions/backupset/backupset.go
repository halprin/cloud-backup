package backupset

import (
	"github.com/halprin/cloud-backup/archival"
	"github.com/halprin/cloud-backup/compression"
	"github.com/halprin/cloud-backup/config"
	"github.com/halprin/cloud-backup/parallel"
	"log"
	"time"
)

func Backup(configFilePath string) error {
	log.Println("Backing-up file set")

	overallConfig, err := config.New(configFilePath)
	if err != nil {
		return err
	}

	nowTime := time.Now()
	overallFolderName := nowTime.Format(time.RFC3339)

	var errorChannels []chan error

	for _, currentBackupFile := range overallConfig.BackupFiles {
		var errorChannel chan error

		func(currentBackupFile config.BackupFileConfiguration, overallConfig config.BackupConfiguration) {
			errorChannel = parallel.InvokeErrorReturnFunction(func() error {
				return backupFile(currentBackupFile, overallConfig, overallFolderName)
			})
		}(currentBackupFile, overallConfig)

		errorChannels = append(errorChannels, errorChannel)
	}

	for _, errorFromBackupFile := range parallel.ConvertChannelsOfErrorToErrorSlice(errorChannels) {
		if errorFromBackupFile != nil {
			return errorFromBackupFile
		}
	}

	log.Println("Backing-up file set complete")
	return nil
}

func backupFile(backupFile config.BackupFileConfiguration, overallConfig config.BackupConfiguration, overallFolderName string) error {
	log.Printf("Backing-up %s (%s)", backupFile.Title, backupFile.Path)

	uploader, err := getDestinationWriter(backupFile, overallConfig, overallFolderName)
	if err != nil {
		return err
	}

	compressor, err := compression.NewCompressor(uploader)
	if err != nil {
		return err
	}

	archiver := archival.NewArchiver(backupFile.Path, compressor.Writer(), backupFile)

	err = archiver.Archive()
	if err != nil {
		log.Printf("Unable to archive %s", backupFile.Title)
		return err
	}

	err = compressor.Close()
	if err != nil {
		log.Printf("Unable finish the compression of %s", backupFile.Title)
		return err
	}

	err = uploader.Close()
	if err != nil {
		log.Printf("Unable to finish the upload of %s", backupFile.Title)
		return err
	}

	log.Printf("Back-up complete for %s", backupFile.Title)
	return nil
}
