package actions

import (
	"bufio"
	"github.com/halprin/cloud-backup-go/archival"
	"github.com/halprin/cloud-backup-go/compression"
	"github.com/halprin/cloud-backup-go/config"
	"github.com/halprin/cloud-backup-go/crypt"
	"github.com/halprin/cloud-backup-go/parallel"
	"log"
	"os"
	"path/filepath"
)

func Backup() error {
	log.Println("Backing-up file set")

	overallConfig, err := config.BackupConfig()
	if err != nil {
		return err
	}

	var errorChannels []chan error

	for _, currentBackupFile := range overallConfig.BackupFiles {
		var errorChannel chan error

		func(currentBackupFile config.BackupFileConfiguration, overallConfig config.BackupConfiguration) {
			errorChannel = parallel.InvokeErrorReturnFunction(func() error {
				return backupFile(currentBackupFile, overallConfig)
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

func backupFile(backupFile config.BackupFileConfiguration, overallConfig config.BackupConfiguration) error {
	log.Printf("Backing-up %s (%s)", backupFile.Title, backupFile.Path)

	outputFile, err := os.Create(filepath.Join(overallConfig.IntermediatePath, backupFile.Title + ".cipher"))
	if err != nil {
		return err
	}
	defer outputFile.Close()

	encryptor := crypt.NewEncryptor(outputFile, overallConfig)

	bufferedWriter := bufio.NewWriterSize(encryptor, 10 * 1024 * 1024)  //buffer in 10 MB increments

	compressor := compression.NewCompressor(bufferedWriter)
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

	err = bufferedWriter.Flush()
	if err != nil {
		log.Printf("Unable to finish the buffering of %s", backupFile.Title)
		return err
	}

	log.Printf("Back-up complete for %s", backupFile.Title)
	return nil
}
