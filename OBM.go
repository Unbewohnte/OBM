package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Unbewohnte/OBM/logger"
	"github.com/Unbewohnte/OBM/manager"
	"github.com/Unbewohnte/OBM/settings"
	"github.com/Unbewohnte/OBM/util"
)

type result struct {
	beatmapName   string
	numberOfDiffs uint
	successful    uint
	failed        uint
}

type job struct {
	beatmap              manager.Beatmap
	replacementImagePath string
	retrievementPath     string
}

func init() {
	exists, err := settings.DoesExist()
	if err != nil {
		logger.LogError(true, err)
	}
	if !exists {
		// settings file does not exist, so create it and exit (assuming that this is the first run)
		settings.Create()
		logger.LogInfo("Successfully created new settings file")
		os.Exit(0)
	}

	logger.LogInfo("Found settings file")
	return
}

func main() {
	var WG sync.WaitGroup

	startingTime := time.Now()

	SETTINGS := settings.Get()

	// creating black image
	if SETTINGS.CreateBlackBGImage.Enabled {
		err := util.CreateBlackBG(SETTINGS.CreateBlackBGImage.Width, SETTINGS.CreateBlackBGImage.Height)
		if err == nil {
			logger.LogInfo("Successfully created black background")
		} else {
			logger.LogWarning(fmt.Sprintf("Could not create black background : %s; continuing to run...", err))
		}
	}

	beatmaps, err := manager.GetBeatmaps(SETTINGS.OsuDir)
	if err != nil {
		logger.LogError(true, "Error getting beatmaps: ", err)
	}
	logger.LogInfo(fmt.Sprintf("Found %d beatmaps", len(beatmaps)))

	// creating jobs for workers
	jobs := make(chan job, len(beatmaps))
	for _, beatmap := range beatmaps {
		jobs <- job{
			beatmap:              beatmap,
			replacementImagePath: SETTINGS.BackgroundReplacement.ReplacementImagePath,
			retrievementPath:     SETTINGS.BackgroundRetrievement.RetrievementPath,
		}
	}
	close(jobs)

	// perform the magic
	results := make(chan result, len(jobs))
	workerPool(jobs, results, SETTINGS.Workers, &WG)
	WG.Wait()
	close(results)

	// extracting results and logging the last message
	var successful, failed uint = 0, 0
	for result := range results {
		successful += result.successful
		failed += result.failed

		logger.LogInfo(fmt.Sprintf("Beatmap: %s; Number of diffs: %d;\n Successful: %d; Failed: %d",
			result.beatmapName, result.numberOfDiffs, result.successful, result.failed))
	}
	total := successful + failed

	logger.LogInfo(fmt.Sprintf("DONE in %v. %d operations successful (%d%%/100%%); %d failed (%d%%/100%%)",
		time.Since(startingTime), successful, successful/total*100, failed, failed/total*100))

}
