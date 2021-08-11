// +build !localDesination

package restore

import (
	"github.com/halprin/cloud-backup-go/config"
	"github.com/halprin/cloud-backup-go/transfer"
	"io"
)

func getSourceReader(overallConfig config.BackupConfiguration, timestamp string, backupFile string) (io.Reader, error) {
	downloader, err := transfer.NewDownloader(overallConfig, timestamp, backupFile)
	if err != nil {
		return nil, err
	}

	downloadReader, err := downloader.Download()
	if err != nil {
		return nil, err
	}

	return downloadReader, nil
}
