package main

import (
	"github.com/halprin/cloud-backup-go/actions"
	"github.com/halprin/cloud-backup-go/crypt"
	"log"
	"os"
)

func main() {
	backupSet()
	//archiveAndCompressAndEncrypt()
	//decrypt()
}

func backupSet() {
	err := actions.Backup()
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
