package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"unbewohnte.xyz/Unbewohnte/OBM/util"
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
		return "", fmt.Errorf("could not process the given path : %s", err)
	}
	if !stat.IsDir() {
		return "", fmt.Errorf("given Osu! directory is not a directory")
	}

	return songsDir, nil
}

// checks for .osu files in given path and returns all found instances
func getDiffs(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("could not read a directory : %s", err)
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
		return nil, fmt.Errorf("could not read a directory : %s", err)
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
