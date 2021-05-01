package manager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Unbewohnte/OBM/logger"
	"github.com/Unbewohnte/OBM/util"
)

// the main beatmap struct, contains necessary data for functions
type Beatmap struct {
	Name  string
	Path  string
	Diffs []string
}

// parses given .osu file and returns the filename of its background
// NOTE: Osu! beatmap (as whole) can have multiple backgrounds for each .osu file
// the perfect example : https://osu.ppy.sh/beatmapsets/43701#osu/137122
// this is why this functions asks for a certain difficulty (.osu filename) to be sure
// to return the correct background name
func (BEATMAP *Beatmap) GetBackgroundName(mapName string) (string, error) {
	beatmapBytes, err := os.ReadFile(filepath.Join(BEATMAP.Path, mapName))
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
		if util.IsImage(chunk) {
			return strings.Split(chunk, "\"")[1], nil
		}
	}
	return "", nil
}

// parses each beatmap`s .osu file for background info;
// removes original background and replaces it with copied version of given image
func (BEATMAP *Beatmap) ReplaceBackgrounds(replacementPicPath string) (successful, failed uint64) {
	// looping through each .osu file of a beatmap
	for _, diff := range BEATMAP.Diffs {
		background, err := BEATMAP.GetBackgroundName(diff)
		if err != nil || background == "" {
			logger.LogError(false, fmt.Sprintf("BEATMAP: %s: Error gettings background filename: %s", diff, err))
			failed++
			continue
		}
		// remove old bg
		err = os.Remove(filepath.Join(BEATMAP.Path, background))
		if err != nil {
			logger.LogError(false, fmt.Sprintf("BEATMAP: %s: Could not remove the old bg: %s", diff, err))
			failed++
			continue
		}

		// copy given picture, thus replacing background
		err = util.CopyFile(replacementPicPath, filepath.Join(BEATMAP.Path, background))
		if err != nil {
			logger.LogError(false, fmt.Sprintf("BEATMAP: %s: Could not copy: %s", diff, err))
			failed++
			continue
		}
		successful++
	}
	return successful, failed
}

// retrieves backgrounds from given beatmap folder (same as in `ReplaceBackgrounds`) and copies them to the retrievement path
func (BEATMAP *Beatmap) RetrieveBackgrounds(retrievementPath string) (successful, failed uint64) {
	// looping through each .osu file of a beatmap
	for _, diff := range BEATMAP.Diffs {
		background, err := BEATMAP.GetBackgroundName(diff)
		if err != nil || background == "" {
			logger.LogError(false, fmt.Sprintf("BEATMAP: %s: Error gettings background filename: %s", diff, err))
			failed++
			continue
		}

		// creating directory that represents current beatmap
		dstPath := filepath.Join(retrievementPath, BEATMAP.Name)

		err = os.MkdirAll(dstPath, os.ModePerm)
		if err != nil {
			logger.LogError(false, fmt.Sprintf("BEATMAP: %s: Error creating a directory: %s", diff, err))
			failed++
			continue
		}

		// copy background to the `retrievementPath`
		err = util.CopyFile(filepath.Join(BEATMAP.Path, background), filepath.Join(dstPath, background))
		if err != nil {
			logger.LogError(false, fmt.Sprintf("BEATMAP: %s: Could not copy: %s", diff, err))
			failed++
			continue
		}
		successful++
	}
	return successful, failed
}
