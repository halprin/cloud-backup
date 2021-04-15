package main

import (
	"bufio"
	"github.com/halprin/cloud-backup-go/archival"
	"github.com/halprin/cloud-backup-go/compression"
	"github.com/halprin/cloud-backup-go/crypt"
	"io/fs"
	"log"
	"os"
)

func main() {
	archiveAndCompressAndEncrypt()
	//encrypt()
	//decrypt()
}

func archiveAndCompress() {
	tarData, err := archival.Archive(os.Args[1])
	if err != nil {
		log.Fatal(err.Error())
	}

	gzipData, err := compression.Compress(tarData)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = os.WriteFile(os.Args[2], gzipData, fs.FileMode(777))
	if err != nil {
		log.Fatal(err.Error())
	}
}

func archiveAndCompressAndEncrypt() {

	outputFile, err := os.Create(os.Args[2])
	if err != nil {
		log.Fatal(err.Error())
	}
	defer outputFile.Close()

	encryptor := crypt.NewEncryptor(outputFile)

	bufferedWriter := bufio.NewWriterSize(encryptor, 1024 * 1024)

	compressor := compression.NewCompressor(bufferedWriter)
	archiver := archival.NewArchiver(os.Args[1], compressor.Writer())

	err = archiver.Archive()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = compressor.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func encrypt() {
	plaintext, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err.Error())
	}

	ciphertext, err := crypt.Encrypt(plaintext)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = os.WriteFile(os.Args[2], ciphertext, 0777)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func decrypt() {
	ciphertext, err := os.ReadFile(os.Args[2])
	if err != nil {
		log.Fatal(err.Error())
	}

	plaintext, err := crypt.Decrypt(ciphertext)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = os.WriteFile(os.Args[1], plaintext, 0777)
	if err != nil {
		log.Fatal(err.Error())
	}
}
