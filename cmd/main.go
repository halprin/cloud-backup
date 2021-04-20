package main

import (
	"github.com/halprin/cloud-backup-go/actions"
	"github.com/halprin/cloud-backup-go/config"
	"github.com/halprin/cloud-backup-go/crypt"
	"log"
	"os"
)

func main() {
	backupSet()
	//decrypt()
}

func backupSet() {
	err := actions.Backup()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func decrypt() {
	overallConfig, err := config.BackupConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	ciphertextFile, err := os.Open(os.Args[2])
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ciphertextFile.Close()

	outputFile, err := os.Create(os.Args[2] + ".tar.gz")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer outputFile.Close()

	decryptor := crypt.NewDecryptor(ciphertextFile, outputFile, overallConfig)

	err = decryptor.Decrypt()
	if err != nil {
		log.Println("decrypt")
		log.Fatal(err.Error())
	}
}
