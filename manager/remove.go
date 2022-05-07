package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"unbewohnte.xyz/Unbewohnte/OBM/logger"
)

// parses each difficulty for background info, removes found backgrounds
func (BEATMAP *Beatmap) RemoveBackgrounds() (successful, failed uint) {
	// looping through each .osu file of a beatmap
	for _, diff := range BEATMAP.Diffs {
		background, err := BEATMAP.GetBackgroundName(diff)
		if err != nil || background == "" {
			logger.LogError(false, fmt.Sprintf("BEATMAP: %s: Error getting background filename: %s", diff, err))
			failed++
			continue
		}
		// remove background
		err = os.Remove(filepath.Join(BEATMAP.Path, background))
		if err != nil {
			// background file does not exist (success ???)
			successful++
			continue
		}
		successful++
	}
	return successful, failed
}
