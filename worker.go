package main

import (
	"sync"
)

// a basic implementation of a concurrent worker
func worker(jobs <-chan job, results chan result, WG *sync.WaitGroup) {
	defer WG.Done()
	for job := range jobs {
		var successful, failed uint = 0, 0

		// the order is: Replace->Retrieve->Remove (if all 3 options are enabled)
		if job.replacementImagePath != "" {
			s, f := job.beatmap.ReplaceBackgrounds(job.replacementImagePath)
			successful += s
			failed += f
		}
		if job.retrievementPath != "" {
			s, f := job.beatmap.RetrieveBackgrounds(job.retrievementPath)
			successful += s
			failed += f
		}
		if job.remove == true {
			s, f := job.beatmap.RemoveBackgrounds()
			successful += s
			failed += f
		}

		results <- result{
			beatmapName:   job.beatmap.Name,
			numberOfDiffs: uint(len(job.beatmap.Diffs)),
			successful:    successful,
			failed:        failed,
		}
	}

}

// the `starter` that `glues` workers and jobs together
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
