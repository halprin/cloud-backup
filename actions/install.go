package actions

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed launchdTemplate.xml
var launchdTemplateString string

var globalDaemons = "/Library/LaunchDaemons/"
var launchdConfigPath = filepath.Join(globalDaemons, "io.halprin.backup.plist")

type templateFields struct {
	Script    string
	ConfigYml string
	Interval  string
}

func Install(configFilePath string, month *int, day *int, weekday *int, hour *int, minute *int) error {
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
