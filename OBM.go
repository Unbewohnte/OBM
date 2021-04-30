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

var (
	WG sync.WaitGroup
)

type job struct {
	beatmapFolderPath    string
	replacementImagePath string
	retrievementPath     string
}

type result struct {
	successful uint64
	failed     uint64
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
	startingTime := time.Now()

	SETTINGS := settings.Get()

	// creating black image
	if SETTINGS.CreateBlackBGImage {
		err := util.CreateBlackBG(1920, 1080)
		if err == nil {
			logger.LogInfo("Successfully created black background")
		} else {
			logger.LogWarning(fmt.Sprintf("Could not create black background : %s; continuing to run...", err))
		}
	}

	beatmaps, err := manager.GetBeatmapFolderPaths(SETTINGS.OsuDir)
	if err != nil {
		logger.LogError(true, "Error getting beatmap folders: ", err)
	}
	logger.LogInfo(fmt.Sprintf("Found %d beatmap folders", len(beatmaps)))

	// creating jobs for workers
	jobs := make(chan job, len(beatmaps))
	for _, beatmap := range beatmaps {
		jobs <- job{
			beatmapFolderPath:    beatmap,
			replacementImagePath: SETTINGS.BackgroundReplacement.ReplacementImagePath,
			retrievementPath:     SETTINGS.BackgroundRetrievement.RetrievementPath,
		}
	}
	close(jobs)

	results := make(chan result, len(jobs))
	workerPool(jobs, results, SETTINGS.Workers, &WG)
	WG.Wait()
	close(results)

	// extracting results and logging the last message
	var successful, failed uint64 = 0, 0
	for result := range results {
		successful += result.successful
		failed += result.failed
	}
	total := successful + failed

	logger.LogInfo(fmt.Sprintf("DONE in %v. %d operations successful (%d%%/100%%); %d failed (%d%%/100%%)",
		time.Since(startingTime), successful, successful/total*100, failed, failed/total*100))

}
