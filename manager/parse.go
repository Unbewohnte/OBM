package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Unbewohnte/OBM/util"
)

// parses given .osu file and returns the filename of its background
// NOTE: Osu! beatmap (as whole) can have multiple backgrounds for each .osu file
// the perfect example : https://osu.ppy.sh/beatmapsets/43701#osu/137122
// this is why this function asks for a certain difficulty (.osu filename) to be sure
// to return the correct background name
func (BEATMAP *Beatmap) GetBackgroundName(diff string) (string, error) {
	beatmapBytes, err := os.ReadFile(filepath.Join(BEATMAP.Path, diff))
	if err != nil {
		return "", err
	}
	beatmapContents := string(beatmapBytes)

	// get index of "[Events]" (this is where BG filename is stored)
	eventsIndex := strings.Index(beatmapContents, "[Events]")
	if eventsIndex == -1 {
		return "", fmt.Errorf("could not retrieve index of \"[Events]\"")
	}
	// get index of [TimingPoints] (this tag is right after the previous "[Events]" tag,
	// so we can grab the whole "[Events]" tag contents)
	timingPointsIndex := strings.Index(beatmapContents, "[TimingPoints]")
	if timingPointsIndex == -1 {
		return "", fmt.Errorf("could not retrieve index of \"[TimingPoints]\"")
	}
	contentBetween := strings.Split(beatmapContents[eventsIndex:timingPointsIndex], ",")

	for _, chunk := range contentBetween {
		if util.IsImage(chunk) {
			return strings.Split(chunk, "\"")[1], nil
		}
	}
	return "", nil
}
