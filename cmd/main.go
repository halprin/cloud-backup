package main

import (
	"github.com/halprin/cloud-backup-go/archival"
	"github.com/halprin/cloud-backup-go/crypt"
	"io/fs"
	"log"
	"os"
)

func main() {
	archive()
	//encrypt()
	//decrypt()
}

func archive() {
	tarFile, err := archival.Archive(os.Args[1])
	if err != nil {
		log.Fatal(err.Error())
	}

	err = os.WriteFile(os.Args[2], tarFile, fs.FileMode(777))
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
