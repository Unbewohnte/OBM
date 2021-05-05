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

// filepath.Joins the main osu directory with its songs folder
func getSongsDir(baseOsuDir string) (string, error) {
	songsDir := filepath.Join(baseOsuDir, "Songs")

	stat, err := os.Stat(songsDir)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could not process the given path : %s", err))
	}
	if !stat.IsDir() {
		return "", errors.New("Given Osu! directory is not a directory !")
	}

	return songsDir, nil
}

// checks for .osu files in given path and returns all found instances
func getDiffs(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not read a directory : %s", err))
	}

	var diffs []string
	for _, file := range files {
		filename := file.Name()
		if util.IsBeatmap(filename) {
			diffs = append(diffs, filename)
		}
	}
	return diffs, nil
}

// constructs a `Beatmap` struct and returns it
func newBeatmap(name, path string, diffs []string) Beatmap {
	return Beatmap{
		Name:  name,
		Path:  path,
		Diffs: diffs,
	}
}

// returns an array of beatmaps from given base Osu! directory
func GetBeatmaps(baseOsuDir string) ([]Beatmap, error) {
	songsDir, err := getSongsDir(baseOsuDir)
	if err != nil {
		return nil, err
	}
	contents, err := os.ReadDir(songsDir)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not read a directory : %s", err))
	}

	var beatmaps []Beatmap
	// looping through all folders in yourOsuDir/Songs/ directory
	for _, file := range contents {
		if file.IsDir() {
			// retrieving all necessary data for creating a new instance of a beatmap
			beatmapName := file.Name()
			pathToBeatmap := filepath.Join(songsDir, beatmapName)
			diffs, err := getDiffs(pathToBeatmap)
			if err != nil {
				continue
			}
			newBeatmap := newBeatmap(beatmapName, pathToBeatmap, diffs)

			beatmaps = append(beatmaps, newBeatmap)
		}
	}

	return beatmaps, nil
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
func (BEATMAP *Beatmap) ReplaceBackgrounds(replacementImgPath string) (successful, failed uint) {
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
		err = util.CopyFile(replacementImgPath, filepath.Join(BEATMAP.Path, background))
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
func (BEATMAP *Beatmap) RetrieveBackgrounds(retrievementPath string) (successful, failed uint) {
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

// Search tries to locate instances of beatmaps with the provided name (or part of the name);
// returns a slice of found beatmaps and a number of searched beatmaps
func Search(beatmaps []Beatmap, name string) ([]Beatmap, uint64) {
	var instances []Beatmap
	var searched uint64 = 0

	// to make a search case-insensitive
	name = strings.ToLower(name)
	for _, beatmap := range beatmaps {
		if strings.Contains(strings.ToLower(beatmap.Name), name) {
			instances = append(instances, beatmap)
		}
		searched++
	}
	return instances, searched
}
