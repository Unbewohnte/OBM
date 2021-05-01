package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/Unbewohnte/OBM/logger"
	"github.com/Unbewohnte/OBM/util"
)

const (
	settingsFilename string = "settings.json"
)

// checks if the settings.json exists in current directory
func DoesExist() (bool, error) {
	files, err := os.ReadDir(".")
	if err != nil {
		return false, errors.New(fmt.Sprintf("Unable to read current directory %s", err))
	}

	for _, file := range files {
		if !file.IsDir() && file.Name() == settingsFilename {
			return true, nil
		}
	}

	return false, nil
}

// creates "settings.json" and sets the flag
func Create() error {
	exists, err := DoesExist()
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	file, err := os.Create(settingsFilename)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to create settings file : %s", err))
	}

	// marshaling default settings
	settingsJson, err := json.MarshalIndent(Settings{
		OsuDir: "",
		BackgroundReplacement: backgroundReplacement{
			Enabled:              true,
			ReplacementImagePath: "",
		},
		BackgroundRetrievement: backgroundRetrievement{
			Enabled:          false,
			RetrievementPath: "",
		},
		CreateBlackBGImage: true,
		Workers:            100,
	}, "", " ")
	if err != nil {
		return errors.New(fmt.Sprintf("Could not marshal settings into file : %s", err))
	}
	file.Write(settingsJson)

	file.Close()

	return nil
}

// unmarshalls settings.json into struct and processes the edge-cases
func Get() Settings {
	settingsFileContents, err := os.ReadFile(settingsFilename)
	if err != nil {
		logger.LogError(true, fmt.Sprintf("Could not read settings file : %s", err))
	}

	var settings Settings
	err = json.Unmarshal(settingsFileContents, &settings)
	if err != nil {
		logger.LogError(true, fmt.Sprintf("Could not unmarshal json file : %s", err))
	}

	// if all features are disabled
	if !settings.BackgroundReplacement.Enabled && !settings.BackgroundRetrievement.Enabled {
		logger.LogInfo("No features enabled. Exiting...")
		os.Exit(0)
	}

	// checking for edge cases or mistakes made in the settings file,
	// enabled and disabled fields
	if settings.BackgroundReplacement.Enabled {
		if settings.BackgroundReplacement.ReplacementImagePath == "" || settings.BackgroundReplacement.ReplacementImagePath == " " {
			logger.LogError(true, "`replacementImagePath` is not specified !")
		} else if !util.IsImage(settings.BackgroundReplacement.ReplacementImagePath) {
			logger.LogError(true, "`replacementImagePath` is pointing to a non-image file !`")
		}
	} else {
		settings.BackgroundReplacement.ReplacementImagePath = ""
	}

	if settings.BackgroundRetrievement.Enabled {
		if settings.BackgroundRetrievement.RetrievementPath == "" || settings.BackgroundRetrievement.RetrievementPath == " " {
			logger.LogError(true, "`retrievementPath` is not specified !")
		}
	} else {
		settings.BackgroundRetrievement.RetrievementPath = ""
	}

	if settings.Workers <= 0 {
		settings.Workers = 1
		logger.LogWarning("`workers` is set to 0 or less. Replaced with 1")
	}

	return settings
}
