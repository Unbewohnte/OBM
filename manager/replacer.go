package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Unbewohnte/OBM/logger"
	"github.com/Unbewohnte/OBM/util"
)

// reads contents of given dir; searches for .osu files; parses them for background info;
// removes original background and replaces it with copied version of given image
func ReplaceBackgrounds(beatmapFolder, replacementPicPath string) (successful, failed uint64) {
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

		// remove old background
		err = os.Remove(pathToBackground)
		if err != nil {
			failed++
			logger.LogWarning(fmt.Sprintf("BEATMAP: %s: Could not remove old background : %s", beatmap, err))
		}

		// create new background file
		bgFile, err := os.Create(pathToBackground)
		if err != nil {
			failed++
			logger.LogWarning(fmt.Sprintf("BEATMAP: %s: Could not create new background file : %s", beatmap, err))
			continue
		}
		bgFile.Close()

		// copy the contents of a given image to the newly created bg file
		err = util.CopyFile(replacementPicPath, pathToBackground)
		if err != nil {
			logger.LogWarning(fmt.Sprintf("BEATMAP: %s: Could not copy file: %s", beatmap, err))
			failed++
			continue
		}
		successful++
	}

	return successful, failed
}
