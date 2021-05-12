package actions

import "log"

func Restore(configFilePath string, timestamp string, backupFile string, restorePath string) error {
	log.Println("Restoring file")
	log.Printf("configFilePath=%s, timestamp=%s, backupFile=%s, restorePath=%s", configFilePath, timestamp, backupFile, restorePath)
	return nil
}
