package manager

import (
	"errors"
	"os"
	"strings"

	"github.com/Unbewohnte/OBM/util"
)

// parses given .osu file and returns the filename of its background
func GetBackgroundName(pathTobeatmap string) (string, error) {
	beatmapBytes, err := os.ReadFile(pathTobeatmap)
	if err != nil {
		return "", err
	}
	beatmapContents := string(beatmapBytes)

	// get index of "[Events]" (this is where BG filename is stored)
	eventsIndex := strings.Index(beatmapContents, "[Events]")
	if eventsIndex == -1 {
		return "", errors.New("Could not retrieve index of \"[Events]\"")
	}
	// get index of [TimingPoints] (this tag is right after the previous "[Events]" tag,
	// so we can grab the whole "[Events]" tag contents)
	timingPointsIndex := strings.Index(beatmapContents, "[TimingPoints]")
	if timingPointsIndex == -1 {
		return "", errors.New("Could not retrieve index of \"[TimingPoints]\"")
	}
	contentBetween := strings.Split(beatmapContents[eventsIndex:timingPointsIndex], ",")

	for _, chunk := range contentBetween {
		if util.IsImage(chunk) {
			return strings.Split(chunk, "\"")[1], nil
		}
	}
	return "", nil
}
