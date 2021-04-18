package actions

import (
	"bufio"
	"github.com/halprin/cloud-backup-go/archival"
	"github.com/halprin/cloud-backup-go/compression"
	"github.com/halprin/cloud-backup-go/config"
	"github.com/halprin/cloud-backup-go/crypt"
	"log"
	"os"
	"path"
)

func Backup() error {
	log.Println("Backing-up file set")

	overallConfig, err := config.BackupConfig()
	if err != nil {
		return err
	}

	for _, currentBackupFile := range overallConfig.BackupFiles {
		err = backupFile(currentBackupFile, overallConfig)
		if err != nil {
			return err
		}
	}

	log.Println("Backing-up file set complete")
	return nil
}

func backupFile(backupFile config.BackupFileConfiguration, overallConfig config.BackupConfiguration) error {
	log.Printf("Backing-up %s (%s)", backupFile.Title, backupFile.Path)

	outputFile, err := os.Create(path.Join(overallConfig.IntermediatePath, backupFile.Title + ".cipher"))
	if err != nil {
		return err
	}
	defer outputFile.Close()

	encryptor := crypt.NewEncryptor(outputFile)

	bufferedWriter := bufio.NewWriterSize(encryptor, 10 * 1024 * 1024)  //buffer in 10 MB increments

	compressor := compression.NewCompressor(bufferedWriter)
	archiver := archival.NewArchiver(backupFile.Path, compressor.Writer())

	err = archiver.Archive()
	if err != nil {
		return err
	}

	err = compressor.Close()
	if err != nil {
		return err
	}

	err = bufferedWriter.Flush()
	if err != nil {
		return err
	}

	log.Printf("Back-up complete for %s", backupFile.Title)
	return nil
}
