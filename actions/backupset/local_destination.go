// +build localDesination

package backupset

import (
	"github.com/halprin/cloud-backup/config"
	"io"
	"os"
	"path/filepath"
)

func getDestinationWriter(backupFile config.BackupFileConfiguration, overallConfig config.BackupConfiguration, overallFolderName string) (io.WriteCloser, error) {
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	fullPath := filepath.Join(homeDirectory, "Desktop", overallFolderName, backupFile.Title + ".cipher")
	err = ensureBaseDirExists(fullPath)
	if err != nil {
		return nil, err
	}

	fileWriter, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}

	return fileWriter, nil
}

func ensureBaseDirExists(thePath string) error {
	baseDir := filepath.Dir(thePath)

	info, err := os.Stat(baseDir)
	if err == nil && info.IsDir() {
		return nil
	}

	return os.MkdirAll(baseDir, 0755)
}
