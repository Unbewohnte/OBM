package manager

import "strings"

// Search tries to locate instances of beatmaps with the provided name (or part of the name);
// returns a slice of found beatmaps and a number of searched beatmaps
func Search(beatmaps []Beatmap, name string) ([]Beatmap, uint64) {
	var instances []Beatmap
	var searched uint64 = 0

	// to make the search case-insensitive
	name = strings.ToLower(name)
	for _, beatmap := range beatmaps {
		if strings.Contains(strings.ToLower(beatmap.Name), name) {
			instances = append(instances, beatmap)
		}
		searched++
	}
	return instances, searched
}
