package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Unbewohnte/OBM/logger"
	"github.com/Unbewohnte/OBM/manager"
	"github.com/Unbewohnte/OBM/settings"
)

var (
	WG sync.WaitGroup
)

type job struct {
	songPath    string
	pathToImage string
}

type result struct {
	successful uint64
	failed     uint64
}

// a basic implementation of a concurrent worker
func worker(jobs <-chan job, results chan result, WG *sync.WaitGroup) {
	defer WG.Done()
	for job := range jobs {
		s, f := manager.ReplaceBackgrounds(job.songPath, job.pathToImage)
		results <- result{
			successful: s,
			failed:     f,
		}
	}

}

func workerPool(jobs chan job, results chan result, numOfWorkers int, WG *sync.WaitGroup) {
	// check if there are less jobs than workers
	if numOfWorkers > len(jobs) {
		numOfWorkers = len(jobs)
	}

	// replacing backgrounds for each beatmap concurrently
	for i := 0; i < numOfWorkers; i++ {
		WG.Add(1)
		go worker(jobs, results, WG)
	}
}

func init() {
	exists, err := settings.CheckSettingsFile()
	if err != nil {
		logger.LogError(true, err)
	}
	if exists {
		logger.LogInfo("Found settings file")
		return
	}

	// settings file does not exist, so create it and exit (assuming that this is the first run)
	settings.CreateSettingsFile()
	os.Exit(0)
}

func main() {
	startingTime := time.Now()

	settings := settings.GetSettings()

	// processing given settings
	if settings.CreateBlackBGImage {
		err := manager.CreateBlackBG(1920, 1080)
		if err == nil {
			logger.LogInfo("Successfully created black background")
		} else {
			logger.LogWarning(fmt.Sprintf("Could not create black background : %s; continuing to run...", err))
		}
	}

	osuSongsDir, err := manager.GetSongsDir(settings.OsuDir)
	if err != nil {
		logger.LogError(true, err)
	}

	if settings.ReplacementImagePath == "" || settings.ReplacementImagePath == " " {
		logger.LogError(true, "Image path not specified ! Specify `pathToimage` in settings file !")
	}

	// reading contents of `Songs` folder
	osuSongsDirContents, err := os.ReadDir(osuSongsDir)
	if err != nil {
		logger.LogError(true, fmt.Sprintf("Error reading osu songs directory : %s", err.Error()))
	}

	// creating jobs for workers
	jobs := make(chan job, len(osuSongsDirContents))
	for _, songDir := range osuSongsDirContents {
		if songDir.IsDir() {
			jobs <- job{
				songPath:    filepath.Join(osuSongsDir, songDir.Name()),
				pathToImage: settings.ReplacementImagePath,
			}
		}
	}
	close(jobs)
	logger.LogInfo(fmt.Sprintf("Found %d song folders", len(jobs)))

	results := make(chan result, len(jobs))
	workerPool(jobs, results, settings.Workers, &WG)

	WG.Wait()
	close(results)

	// extracting results and logging the last message
	var successful, failed uint64 = 0, 0
	for result := range results {
		successful += result.successful
		failed += result.failed
	}
	total := successful + failed

	logger.LogInfo(fmt.Sprintf("DONE in %v. %d successful (%d%%/100%%); %d failed (%d%%/100%%)",
		time.Since(startingTime), successful, successful/total*100, failed, failed/total*100))

}
