package main

import (
	"sync"

	"github.com/Unbewohnte/OBM/manager"
)

// a basic implementation of a concurrent worker
func worker(jobs <-chan job, results chan result, WG *sync.WaitGroup) {
	defer WG.Done()
	for job := range jobs {
		var successful, failed uint64 = 0, 0

		if job.retrievementPath != "" {
			s, f := manager.RetrieveBackgrounds(job.beatmapFolderPath, job.retrievementPath)
			successful += s
			failed += f
		}
		if job.replacementImagePath != "" {
			s, f := manager.ReplaceBackgrounds(job.beatmapFolderPath, job.replacementImagePath)
			successful += s
			failed += f
		}
		results <- result{
			successful: successful,
			failed:     failed,
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
