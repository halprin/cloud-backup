package actions

import (
	"log"
	"os"
	"os/exec"
)

func Uninstall() error {
	log.Println("Uninstalling launchd daemon")

	err := unloadLaunchdConfig()
	if err != nil {
		return err
	}

	err = removeLaunchdConfig()
	if err != nil {
		return err
	}

	log.Println("Uninstalled launchd daemon")

	return nil
}

func unloadLaunchdConfig() error {
	log.Println("Unloading from launchd")

	command := exec.Command("launchctl", "unload", launchdConfigPath)
	err := command.Run()
	return err
}

func removeLaunchdConfig() error {
	log.Printf("Deleting the launchd configuration at %s", launchdConfigPath)

	err := os.Remove(launchdConfigPath)
	if err != nil {
		return err
	}

	return nil
}