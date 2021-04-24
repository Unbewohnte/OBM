package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/Unbewohnte/OBM/logger"
)

// struct for json settings` file contents
type Settings struct {
	OsuDir               string `json:"pathToOsu"`
	ReplacementImagePath string `json:"pathToimage"`
	CreateBlackBGImage   bool   `json:"createBlackBackgoundImage"`
	Workers              int    `json:"Workers"`
}

const (
	settingsFilename string = "settings.json"
)

var (
	defaultSettings Settings = Settings{
		OsuDir:               "",
		ReplacementImagePath: "",
		CreateBlackBGImage:   true,
		Workers:              100,
	}
)

// checks if the settings.json exists in current directory
func CheckSettingsFile() (bool, error) {
	files, err := os.ReadDir(".")
	if err != nil {
		return false, errors.New(fmt.Sprintf("ERROR : Unable to read current directory %s", err))
	}

	for _, file := range files {
		if !file.IsDir() && file.Name() == settingsFilename {
			return true, nil
		}
	}

	return false, nil
}

// creates "settings.json" and sets the flag
func CreateSettingsFile() error {
	exists, err := CheckSettingsFile()
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	file, err := os.Create(settingsFilename)
	if err != nil {
		return errors.New(fmt.Sprintf("ERROR: Error creating settings file... : %s", err))
	}

	settingsJson, err := json.MarshalIndent(defaultSettings, "", " ")
	if err != nil {
		return errors.New(fmt.Sprintf("ERROR: Error creating settings file... : %s", err))
	}
	file.Write(settingsJson)

	file.Close()
	logger.LogInfo("Successfully created new settingsFile")

	return nil
}

// unmarshalls settings.json into struct
func GetSettings() Settings {
	settingsFile, err := os.ReadFile(settingsFilename)
	if err != nil {
		logger.LogError(true, "Could not read settings file : ", err.Error())
	}

	var settings Settings
	err = json.Unmarshal(settingsFile, &settings)
	if err != nil {
		logger.LogError(true, "Could not unmarshal json file : ", err)
	}

	if settings.Workers <= 0 {
		logger.LogInfo("`Workers` is set to 0 or less. Replaced with 1")
		settings.Workers = 1
	}

	return settings
}
