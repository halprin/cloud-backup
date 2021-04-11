package main

import (
	"github.com/halprin/cloud-backup-go/crypt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	plaintext, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err.Error())
	}

	cyphertext := crypt.Encrypt(plaintext, "context string")

	err = ioutil.WriteFile(os.Args[2], cyphertext, 0777)
	if err != nil {
		log.Fatal(err.Error())
	}
}
