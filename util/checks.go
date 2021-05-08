package util

import (
	"os"
	"strings"
)

// checks if given path is a directory
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		// path error
		return false
	}
	if !info.IsDir() {
		return false
	}
	return true
}

// checks if given string contains ".osu" file extention
func IsBeatmap(filename string) bool {
	if len(filename) < 5 {
		// too short filename to be a beatmap file
		return false
	}
	if filename[len(filename)-4:] == ".osu" {
		return true
	}
	return false
}

// checks if given string contains the image file extention
func IsImage(filename string) bool {
	if IsDir(filename) {
		// given filename is actually a directory
		return false
	}

	var imageExtentions []string = []string{"jpeg", "jpg", "png", "JPEG", "JPG", "PNG"}
	for _, extention := range imageExtentions {
		if strings.Contains(filename, extention) {
			return true
		}
	}
	return false
}

// checks if given directory/file does exist
func DoesExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}
