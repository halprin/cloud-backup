package main

import (
	"github.com/halprin/cloud-backup-go/config"
	"github.com/halprin/cloud-backup-go/crypt"
	"github.com/halprin/cloud-backup-go/external/cli"
	"log"
	"os"
)

func main() {
	cli.Cli()
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
		log.Println("decrypt error")
		log.Fatal(err.Error())
	}
}
