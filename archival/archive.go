package archival

import (
	"archive/tar"
	"bytes"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func Archive(filePath string) ([]byte, error) {
	byteBuffer := bytes.Buffer{}
	tarWriter := tar.NewWriter(&byteBuffer)
	defer tarWriter.Close()

	err := filepath.WalkDir(filePath, func(currentPath string, fileMetadata fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		err = writeFileMetadataToTar(currentPath, fileMetadata, tarWriter)
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

func writeFileMetadataToTar(currentPath string, fileMetadata fs.DirEntry, tarWriter *tar.Writer) error {
	fileInfo, err := fileMetadata.Info()
	if err != nil {
		return err
	}

	tarHeader, err := tar.FileInfoHeader(fileInfo, currentPath)
	if err != nil {
		return err
	}
	tarHeader.Name = currentPath //because fs.FileInfo's Name method only returns the base name

	err = tarWriter.WriteHeader(tarHeader)
	return err
}
