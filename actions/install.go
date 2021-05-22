package actions

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed launchdTemplate.xml
var launchdTemplateString string

var globalDaemons = "/Library/LaunchDaemons/"
var launchdId = "io.halprin.backup"
var launchdConfigPath = filepath.Join(globalDaemons, fmt.Sprintf("%s.plist", launchdId))

type templateFields struct {
	Script    string
	ConfigYml string
	Interval  string
}

func Install(configFilePath string, month *int, day *int, weekday *int, hour *int, minute *int) error {
	log.Println("Installing launchd daemon")

	if isDaemonLoaded() {
		err := Uninstall()
		if err != nil {
			return err
		}
	}

	err := writeOutLaunchdConfig(configFilePath, month, day, weekday, hour, minute)
	if err != nil {
		return err
	}

	err = loadLaunchdConfig()
	if err != nil {
		log.Println("Error loading launchd config")
		_ = removeLaunchdConfig()
		return err
	}

	log.Println("Launchd daemon installed")
	return nil
}

func loadLaunchdConfig() error {
	log.Println("Loading into launchd")

	command := exec.Command("launchctl", "load", launchdConfigPath)
	err := command.Run()
	return err
}

func writeOutLaunchdConfig(configFilePath string, month *int, day *int, weekday *int, hour *int, minute *int) error {
	log.Printf("Writing out launchd daemon to %s", configFilePath)

	launchdTemplate, err := template.New("launchdTemplate").Parse(launchdTemplateString)
	if err != nil {
		return err
	}

	interval := constructInterval(month, day, weekday, hour, minute)
	fields := templateFields{
		Script:    os.Args[0],
		ConfigYml: configFilePath,
		Interval:  interval,
	}

	endAgentString := &strings.Builder{}
	err = launchdTemplate.Execute(endAgentString, fields)
	if err != nil {
		return err
	}

	err = os.WriteFile(launchdConfigPath, []byte(endAgentString.String()),0644)
	return err
}

func constructInterval(month *int, day *int, weekday *int, hour *int, minute *int) string {
	intervalBuilder := &strings.Builder{}

	if month != nil {
		intervalBuilder.WriteString("<key>Month</key>")
		intervalBuilder.WriteString(fmt.Sprintf("<integer>%d</integer>", *month))
	}

	if day != nil {
		intervalBuilder.WriteString("<key>Day</key>")
		intervalBuilder.WriteString(fmt.Sprintf("<integer>%d</integer>", *day))
	}

	if weekday != nil {
		intervalBuilder.WriteString("<key>Weekday</key>")
		intervalBuilder.WriteString(fmt.Sprintf("<integer>%d</integer>", *weekday))
	}

	if hour != nil {
		intervalBuilder.WriteString("<key>Hour</key>")
		intervalBuilder.WriteString(fmt.Sprintf("<integer>%d</integer>", *hour))
	}

	if minute != nil {
		intervalBuilder.WriteString("<key>Minute</key>")
		intervalBuilder.WriteString(fmt.Sprintf("<integer>%d</integer>", *minute))
	}

	return intervalBuilder.String()
}

func isDaemonLoaded() bool {
	command := exec.Command("launchctl", "list", launchdId)
	err := command.Run()
	if err != nil {
		return false
	}

	return true
}
