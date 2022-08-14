package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"unbewohnte/OBM/logger"
	"unbewohnte/OBM/util"
)

// parses each beatmap`s .osu file for background info;
// removes original background and replaces it with copied version of given image
func (BEATMAP *Beatmap) ReplaceBackgrounds(replacementImgPath string) (successful, failed uint) {
	// looping through each .osu file of a beatmap
	for _, diff := range BEATMAP.Diffs {
		background, err := BEATMAP.GetBackgroundName(diff)
		if err != nil || background == "" {
			logger.LogError(false, fmt.Sprintf("BEATMAP: %s: Error getting background filename: %s", diff, err))
			failed++
			continue
		}
		// remove old bg (if there is no background file - no need to worry)
		os.Remove(filepath.Join(BEATMAP.Path, background))

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
