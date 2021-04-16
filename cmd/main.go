package main

import (
	"bufio"
	"github.com/halprin/cloud-backup-go/archival"
	"github.com/halprin/cloud-backup-go/compression"
	"github.com/halprin/cloud-backup-go/crypt"
	"log"
	"os"
)

func main() {
	archiveAndCompressAndEncrypt()
	decrypt()
}

func archiveAndCompressAndEncrypt() {

	outputFile, err := os.Create(os.Args[2])
	if err != nil {
		log.Fatal(err.Error())
	}
	defer outputFile.Close()

	encryptor := crypt.NewEncryptor(outputFile)

	bufferedWriter := bufio.NewWriterSize(encryptor, 10 * 1024 * 1024)  //buffer in 10 MB increments

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

	err = bufferedWriter.Flush()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func decrypt() {
	ciphertext, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err.Error())
	}

	plaintext, err := crypt.Decrypt(ciphertext)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = os.WriteFile(os.Args[2], plaintext, 0777)
	if err != nil {
		log.Fatal(err.Error())
	}
}
