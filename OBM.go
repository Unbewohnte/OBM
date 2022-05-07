package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"unbewohnte.xyz/Unbewohnte/OBM/logger"
	"unbewohnte.xyz/Unbewohnte/OBM/manager"
	"unbewohnte.xyz/Unbewohnte/OBM/settings"
	"unbewohnte.xyz/Unbewohnte/OBM/util"
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
	remove               bool
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

	// parse for `-beatmap` argument
	flag.Parse()

	// if `-showOrder` is checked - show the message and exit
	if *showOrder == true {
		fmt.Print(orderMessage)
		os.Exit(0)
	}
	return
}

func main() {
	var WG sync.WaitGroup

	startingTime := time.Now()

	SETTINGS := settings.Get()

	// creating black image if enabled
	if SETTINGS.CreateBlackBGImage.Enabled {
		err := util.CreateBlackBG(SETTINGS.CreateBlackBGImage.Width, SETTINGS.CreateBlackBGImage.Height)
		if err == nil {
			logger.LogInfo("Successfully created black background")
		} else {
			logger.LogWarning(fmt.Sprintf("Could not create black background : %s; continuing to run...", err))
		}
	}

	// get an array of all beatmaps
	beatmaps, err := manager.GetBeatmaps(SETTINGS.OsuDir)
	if err != nil {
		logger.LogError(true, "Error getting beatmaps: ", err)
	}
	logger.LogInfo(fmt.Sprintf("Found %d beatmaps", len(beatmaps)))

	// If `-beatmap` flag is specified - do the magic only on found beatmaps
	if *cmdlnBeatmap != "" {
		logger.LogInfo(fmt.Sprintf("Trying to locate \"%s\"...", *cmdlnBeatmap))
		found, n := manager.Search(beatmaps, *cmdlnBeatmap)
		logger.LogInfo(fmt.Sprintf("Checked %d beatmaps. Found %d instance(s)", n, len(found)))

		// if found nothing - exit
		if len(found) == 0 {
			os.Exit(0)
		}
		// replace all beatmaps with found ones
		beatmaps = found
	}

	// creating jobs for workers
	jobs := make(chan job, len(beatmaps))
	for _, beatmap := range beatmaps {
		jobs <- job{
			beatmap:              beatmap,
			replacementImagePath: SETTINGS.BackgroundReplacement.ReplacementImagePath,
			retrievementPath:     SETTINGS.BackgroundRetrievement.RetrievementPath,
			remove:               SETTINGS.BackgroundRemovement.Enabled,
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

	logger.LogInfo(fmt.Sprintf("DONE in %v. %d operations successful (%.2f%%/100%%); %d failed (%.2f%%/100%%)",
		time.Since(startingTime), successful, float32(successful)/float32(total)*100, failed, float32(failed)/float32(total)*100))

}
