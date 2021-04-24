package manager

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// checks if given string contains ".osu" file extention
func isBeatmap(filename string) bool {
	if len(filename) < 5 {
		return false
	}
	if filename[len(filename)-4:] == ".osu" {
		return true
	}
	return false
}

// checks if given string contains the image file extention
func isImage(filename string) bool {
	var imageExtentions []string = []string{"jpeg", "jpg", "png", "JPEG", "JPG", "PNG"}
	for _, extention := range imageExtentions {
		if strings.Contains(filename, extention) {
			return true
		}
	}
	return false
}

// opens given files, copies one into another
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not open src file : %s", err))
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_WRONLY, os.ModePerm)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not open dst file : %s", err))
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not copy file : %s", err))
	}

	return nil
}
