package archival

import (
	"archive/tar"
	"github.com/halprin/cloud-backup-go/config"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type archiver struct {
	filePath   string
	tarWriter  *tar.Writer
	fileConfig config.BackupFileConfiguration
}

func NewArchiver(filePath string, outputWriter io.Writer, fileConfig config.BackupFileConfiguration) *archiver {
	tarWriter := tar.NewWriter(outputWriter)
	return &archiver{
		filePath: filePath,
		tarWriter: tarWriter,
		fileConfig: fileConfig,
	}
}

func (receiver *archiver) Archive() error {
	parentDirectoryPath := filepath.Dir(receiver.filePath)

	defer receiver.tarWriter.Close()

	err := filepath.WalkDir(receiver.filePath, func(currentPath string, fileMetadata fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if receiver.shouldSkipArchival(currentPath) {
			if fileMetadata.IsDir() {
				return filepath.SkipDir
			} else {
				return nil
			}
		}

		//get the relative path the tar file doesn't have the entire absolute path but just the relative path from the start of the walk down to the current file
		relativePath, err := filepath.Rel(parentDirectoryPath, currentPath)
		if err != nil {
			return err
		}

		err = receiver.writeFileMetadataToTar(relativePath, fileMetadata)
		if err != nil {
			return err
		}

		if !fileMetadata.Type().IsRegular() {
			//we're done if this is just a directory, symlink, or some non-regular file
			return nil
		}

		//it's an actual file
		err = receiver.writeFileToTar(currentPath)
		return err
	})

	if err != nil {
		return err
	}

	err = receiver.tarWriter.Close()
	return err
}

func (receiver *archiver) Writer() io.Writer {
	return receiver.tarWriter
}

func (receiver *archiver) shouldSkipArchival(filePath string) bool {
	matches := false

	for _, currentIgnorePath := range receiver.fileConfig.Ignore {
		if strings.Contains(filePath, currentIgnorePath) {
			matches = true
			break
		}
	}

	return matches
}

func (receiver *archiver) writeFileMetadataToTar(relativePath string, fileMetadata fs.DirEntry) error {
	fileInfo, err := fileMetadata.Info()
	if err != nil {
		return err
	}

	tarHeader, err := tar.FileInfoHeader(fileInfo, relativePath)
	if err != nil {
		return err
	}
	tarHeader.Name = relativePath //because fs.FileInfo's Name method only returns the base name

	err = receiver.tarWriter.WriteHeader(tarHeader)
	return err
}

func (receiver *archiver) writeFileToTar(currentPath string) error {
	fileData, err := os.Open(currentPath)
	if err != nil {
		return err
	}
	defer fileData.Close()

	_, err = io.Copy(receiver.tarWriter, fileData)
	return err
}
