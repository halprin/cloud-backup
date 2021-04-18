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
	"sync"
)

func Backup() error {
	log.Println("Backing-up file set")

	overallConfig, err := config.BackupConfig()
	if err != nil {
		return err
	}

	waitGroup := &sync.WaitGroup{}
	for _, currentBackupFile := range overallConfig.BackupFiles {
		waitGroup.Add(1)
		go backupFile(currentBackupFile, overallConfig, waitGroup)
	}

	waitGroup.Wait()

	log.Println("Backing-up file set complete")
	return nil
}

func backupFile(backupFile config.BackupFileConfiguration, overallConfig config.BackupConfiguration, waitGroup *sync.WaitGroup) error {
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
	waitGroup.Done()
	return nil
}
