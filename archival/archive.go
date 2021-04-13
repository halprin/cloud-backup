package archival

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func Archive(filePath string) ([]byte, error) {
	parentDirectoryPath := filepath.Dir(filePath)

	byteBuffer := bytes.Buffer{}
	tarWriter := tar.NewWriter(&byteBuffer)
	defer tarWriter.Close()

	err := filepath.WalkDir(filePath, func(currentPath string, fileMetadata fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		//get the relative path the tar file doesn't have the entire absolute path but just the relative path from the start of the walk down to the current file
		relativePath, err := filepath.Rel(parentDirectoryPath, currentPath)
		if err != nil {
			return err
		}

		fmt.Println(relativePath)

		err = writeFileMetadataToTar(relativePath, fileMetadata, tarWriter)
		if err != nil {
			return err
		}

		if fileMetadata.IsDir() {
			//we're done if this is just a directory
			return nil
		}

		//it's an actual file
		err = writeFileToTar(currentPath, tarWriter)
		return err
	})

	if err != nil {
		return nil, err
	}

	return byteBuffer.Bytes(), nil
}

func writeFileToTar(currentPath string, tarWriter *tar.Writer) error {
	fileData, err := os.Open(currentPath)
	if err != nil {
		return err
	}
	defer fileData.Close()

	_, err = io.Copy(tarWriter, fileData)
	return err
}

func writeFileMetadataToTar(relativePath string, fileMetadata fs.DirEntry, tarWriter *tar.Writer) error {
	fileInfo, err := fileMetadata.Info()
	if err != nil {
		return err
	}

	tarHeader, err := tar.FileInfoHeader(fileInfo, relativePath)
	if err != nil {
		return err
	}
	tarHeader.Name = relativePath //because fs.FileInfo's Name method only returns the base name

	err = tarWriter.WriteHeader(tarHeader)
	return err
}
