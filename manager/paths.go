package manager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

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

// returns an array of full filepaths to each beatmap from given base Osu! directory
func GetBeatmapFolderPaths(baseOsuDir string) ([]string, error) {
	songsDir, err := getSongsDir(baseOsuDir)
	if err != nil {
		return nil, err
	}
	contents, err := os.ReadDir(songsDir)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not read a directory : %s", err))
	}

	var beatmapFolderPaths []string
	for _, file := range contents {
		if file.IsDir() {
			path := filepath.Join(songsDir, file.Name())
			beatmapFolderPaths = append(beatmapFolderPaths, path)
		}
	}

	return beatmapFolderPaths, nil
}
