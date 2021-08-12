// +build localDesination

package restore

import (
	"github.com/halprin/cloud-backup/config"
	"io"
	"os"
	"path/filepath"
)

func getSourceReader(overallConfig config.BackupConfiguration, timestamp string, backupFile string) (io.Reader, error) {
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	fullPath := filepath.Join(homeDirectory, "Desktop", timestamp, backupFile + ".cipher")

	fileReader, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}

	return fileReader, nil
}
