package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Unbewohnte/OBM/logger"
	"github.com/Unbewohnte/OBM/manager"
	"github.com/Unbewohnte/OBM/settings"
)

var (
	WG sync.WaitGroup
)

// creates a complete black image file
func createBlackBG(width, height int) error {
	bgfile, err := os.Create("blackBG.png")
	if err != nil {
		return errors.New(fmt.Sprintf("Could not create black background file : %s", err))
	}
	image := image.NewRGBA(image.Rect(0, 0, width, height))
	bounds := image.Bounds()

	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			image.Set(x, y, color.Black)
		}
	}
	err = png.Encode(bgfile, image)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not encode an image : %s", err))
	}
	err = bgfile.Close()
	if err != nil {
		return errors.New(fmt.Sprintf("Could not close the background file : %s", err))
	}

	return nil
}

// a basic implementation of a concurrent worker
func worker(paths <-chan string, replacementImage string, successful, failed *uint64, WG *sync.WaitGroup) {
	defer WG.Done()
	for songPath := range paths {
		s, f := manager.ReplaceBackgrounds(songPath, replacementImage)
		*successful += s
		*failed += f
	}

}

func init() {
	exists, err := settings.CheckSettingsFile()
	if err != nil {
		logger.LogError(true, err)
	}
	if exists {
		logger.LogInfo("Found settings file")
		return
	}

	// settings file does not exist, so create it and exit
	settings.CreateSettingsFile()
	os.Exit(0)
}

func main() {
	startingTime := time.Now().UTC()

	settings := settings.GetSettings()

	// process the given settings
	if settings.CreateBlackBGImage {
		err := createBlackBG(1920, 1080)
		if err == nil {
			logger.LogInfo("Successfully created black background")
		} else {
			logger.LogWarning(fmt.Sprintf("Could not create black background : %s; continuing to run...", err))
		}
	}

	osuSongsDir, err := manager.GetSongsDir(settings.OsuDir)
	if err != nil {
		logger.LogError(true, err)
	}

	if settings.ReplacementImagePath == "" || settings.ReplacementImagePath == " " {
		logger.LogError(true, "Image path not specified ! Specify `pathToimage` in settings file !")
	}

	// reading contents of `Songs` folder
	osuSongsDirContents, err := os.ReadDir(osuSongsDir)
	if err != nil {
		logger.LogError(true, fmt.Sprintf("Error reading osu songs directory : %s", err.Error()))
	}

	// storing all paths to each beatmap
	songPaths := make(chan string, len(osuSongsDirContents))
	for _, songDir := range osuSongsDirContents {
		if songDir.IsDir() {
			songPaths <- filepath.Join(osuSongsDir, songDir.Name())
		}
	}
	logger.LogInfo(fmt.Sprintf("Found %d song folders", len(songPaths)))

	// check if there is less job than workers
	if settings.Workers > len(songPaths) {
		settings.Workers = len(songPaths)
	}

	// replacing backgrounds for each beatmap concurrently
	var successful, failed uint64 = 0, 0
	for i := 0; i < int(settings.Workers); i++ {
		WG.Add(1)
		go worker(songPaths, settings.ReplacementImagePath, &successful, &failed, &WG)
	}

	close(songPaths)
	WG.Wait()

	endTime := time.Now().UTC()

	logger.LogInfo(fmt.Sprintf("\n\nDONE in %v . %d successful; %d failed", endTime.Sub(startingTime), successful, failed))

}
