package manager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Unbewohnte/OBM/logger"
)

// filepath.Joins the main osu directory with its songs folder
func GetSongsDir(osudir string) (string, error) {
	songsDir := filepath.Join(osudir, "Songs")

	stat, err := os.Stat(songsDir)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could not read the given path : %s", err))
	}
	if !stat.IsDir() {
		return "", errors.New("Given Osu! directory is not a directory !")
	}

	return songsDir, nil
}

// parses given .osu file and returns the filename of its background
func GetBackgroundName(pathToOSUbeatmap string) (string, error) {
	beatmapBytes, err := os.ReadFile(pathToOSUbeatmap)
	if err != nil {
		return "", err
	}
	beatmapContents := string(beatmapBytes)

	// get index of "[Events]" (this is where BG filename is stored)
	eventsIndex := strings.Index(beatmapContents, "[Events]")
	if eventsIndex == -1 {
		return "", errors.New("Could not retrieve index of \"[Events]\"")
	}
	// get index of [TimingPoints] (this tag is right after the previous "[Events]" tag,
	// so we can grab the whole "[Events]" tag contents)
	timingPointsIndex := strings.Index(beatmapContents, "[TimingPoints]")
	if timingPointsIndex == -1 {
		return "", errors.New("Could not retrieve index of \"[TimingPoints]\"")
	}
	contentBetween := strings.Split(beatmapContents[eventsIndex:timingPointsIndex], ",")

	for _, chunk := range contentBetween {
		if isImage(chunk) {
			return strings.Split(chunk, "\"")[1], nil
		}
	}
	return "", nil
}

// reads contents of given dir; searches for .osu files; parses them for background info;
// removes original background and replaces it with copied version of given image
func ReplaceBackgrounds(beatmapFolder, replacementPicPath string) (successful, failed uint64) {
	files, err := os.ReadDir(beatmapFolder)
	if err != nil {
		logger.LogError(false, fmt.Sprintf("Wrong path : %s", err))
	}
	for _, file := range files {
		filename := file.Name()

		if isBeatmap(filename) {
			beatmap := filename

			// getting BG filename
			beatmapBackgroundFilename, err := GetBackgroundName(filepath.Join(beatmapFolder, beatmap))
			if err != nil {
				logger.LogWarning(fmt.Sprintf("BEATMAP: %s: Error getting background filename: %s", beatmap, err))
				failed++
				continue
			}
			if beatmapBackgroundFilename == "" {
				logger.LogWarning(fmt.Sprintf("BEATMAP: %s Could not find background filename in this beatmap file", beatmap))
				failed++
				continue
			}

			backgroundPath := filepath.Join(beatmapFolder, beatmapBackgroundFilename)

			// remove old background
			err = os.Remove(backgroundPath)
			if err != nil {
				failed++
				logger.LogWarning(fmt.Sprintf("BEATMAP: %s: Could not remove old background : %s", beatmap, err))
			}

			// create new background file
			bgFile, err := os.Create(backgroundPath)
			if err != nil {
				failed++
				logger.LogWarning(fmt.Sprintf("BEATMAP: %s: Could not create new background file : %s", beatmap, err))
				continue
			}
			defer bgFile.Close()

			// copy the contents of a given image to the newly created bg file
			err = copyFile(replacementPicPath, backgroundPath)
			if err != nil {
				logger.LogWarning(fmt.Sprintf("BEATMAP: %s: Could not copy file: %s", beatmap, err))
				failed++
				continue
			}
			successful++
		}

	}
	return successful, failed
}
