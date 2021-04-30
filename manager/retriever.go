package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Unbewohnte/OBM/logger"
	"github.com/Unbewohnte/OBM/util"
)

// retrieves backgrounds from given beatmap folder (same as in `ReplaceBackgrounds`) and copies them to the retrievement path
func RetrieveBackgrounds(beatmapFolder, retrievementPath string) (successful, failed uint64) {
	files, err := os.ReadDir(beatmapFolder)
	if err != nil {
		logger.LogError(true, fmt.Sprintf("Could not read directory : %s", err))
	}
	for _, file := range files {
		filename := file.Name()

		// if not a beatmap - skip
		if !util.IsBeatmap(filename) {
			continue
		}

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

		pathToBackground := filepath.Join(beatmapFolder, beatmapBackgroundFilename)

		// creating a directory with the name of current beatmap folder in the retrievement path
		dstPath := filepath.Join(retrievementPath, filepath.Base(beatmapFolder))

		err = os.MkdirAll(dstPath, os.ModePerm)
		if err != nil {
			logger.LogWarning(fmt.Sprintf("BEATMAP: %s: Could not create a directory (%s) for copying", beatmap, dstPath))
			continue
		}

		// creating a copy file
		fullPathToCopy := filepath.Join(dstPath, beatmapBackgroundFilename)
		dstFile, err := os.Create(fullPathToCopy)
		if err != nil {
			logger.LogWarning(fmt.Sprintf("BEATMAP: %s: Could not create a copy file", beatmap))
			failed++
			continue
		}
		dstFile.Close()

		// copy the background file to the retrievement path
		err = util.CopyFile(pathToBackground, fullPathToCopy)
		if err != nil {
			logger.LogWarning(fmt.Sprintf("BEATMAP: %s: Could not copy file: %s", beatmap, err))
			failed++
			continue
		}
		successful++
	}

	return successful, failed
}
