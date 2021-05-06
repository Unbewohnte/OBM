package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Unbewohnte/OBM/logger"
	"github.com/Unbewohnte/OBM/util"
)

// retrieves backgrounds from given beatmap folder (same as in `ReplaceBackgrounds`) and copies them to the retrievement path
func (BEATMAP *Beatmap) RetrieveBackgrounds(retrievementPath string) (successful, failed uint) {
	// looping through each .osu file of a beatmap
	for _, diff := range BEATMAP.Diffs {
		background, err := BEATMAP.GetBackgroundName(diff)
		if err != nil || background == "" {
			logger.LogError(false, fmt.Sprintf("BEATMAP: %s: Error getting background filename: %s", diff, err))
			failed++
			continue
		}

		// check if the background does exist
		exists := util.DoesExist(filepath.Join(BEATMAP.Path, background))
		if !exists {
			// if not - we cannot copy it, so moving to the next diff
			logger.LogWarning(fmt.Sprintf("BEATMAP: %s: Background does not exist", diff))
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
