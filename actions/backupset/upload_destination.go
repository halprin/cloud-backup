//go:build !localDesination

package backupset

import (
	"github.com/halprin/cloud-backup/config"
	"github.com/halprin/cloud-backup/transfer"
	"io"
)

func getDestinationWriter(backupFile config.BackupFileConfiguration, overallConfig config.BackupConfiguration, overallFolderName string) (io.WriteCloser, error) {
	uploader, err := transfer.NewUploader(backupFile, overallConfig, overallFolderName)
	if err != nil {
		return nil, err
	}

	return uploader, nil
}
