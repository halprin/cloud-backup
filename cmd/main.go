package main

import (
	"github.com/halprin/cloud-backup-go/crypt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	//encrypt()
	decrypt()
}

func encrypt() {
	plaintext, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err.Error())
	}

	ciphertext, err := crypt.Encrypt(plaintext)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = ioutil.WriteFile(os.Args[2], ciphertext, 0777)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func decrypt() {
	ciphertext, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		log.Fatal(err.Error())
	}

	plaintext, err := crypt.Decrypt(ciphertext)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = ioutil.WriteFile(os.Args[1], plaintext, 0777)
	if err != nil {
		log.Fatal(err.Error())
	}
}
