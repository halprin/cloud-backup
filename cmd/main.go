package main

import (
	"bufio"
	"github.com/halprin/cloud-backup-go/archival"
	"github.com/halprin/cloud-backup-go/compression"
	"github.com/halprin/cloud-backup-go/config"
	"github.com/halprin/cloud-backup-go/crypt"
	"log"
	"os"
)

func main() {
	//archiveAndCompressAndEncrypt()
	//decrypt()
	configuration()
}

func configuration() {
	theConfig, err := config.BackupConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println(theConfig.S3Bucket)
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
	ciphertextFile, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ciphertextFile.Close()

	outputFile, err := os.Create(os.Args[2])
	if err != nil {
		log.Fatal(err.Error())
	}
	defer outputFile.Close()

	decryptor := crypt.NewDecryptor(ciphertextFile, outputFile)

	err = decryptor.Decrypt()
	if err != nil {
		log.Println("decrypt")
		log.Fatal(err.Error())
	}
}
