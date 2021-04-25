package manager

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
	"strings"
)

// creates a complete black image file
func CreateBlackBG(width, height int) error {
	bgfile, err := os.Create("blackBG.png")
	if err != nil {
		return errors.New(fmt.Sprintf("Could not create black background file : %s", err))
	}
	image := image.NewRGBA(image.Rect(0, 0, width, height))
	bounds := image.Bounds()

	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			image.Set(x, y, color.Black)
		}
	}
	err = png.Encode(bgfile, image)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not encode an image : %s", err))
	}
	err = bgfile.Close()
	if err != nil {
		return errors.New(fmt.Sprintf("Could not close the background file : %s", err))
	}

	return nil
}

// checks if given string contains ".osu" file extention (NOT EXPORTED !)
func isBeatmap(filename string) bool {
	if len(filename) < 5 {
		return false
	}
	if filename[len(filename)-4:] == ".osu" {
		return true
	}
	return false
}

// checks if given string contains the image file extention (NOT EXPORTED !)
func isImage(filename string) bool {
	var imageExtentions []string = []string{"jpeg", "jpg", "png", "JPEG", "JPG", "PNG"}
	for _, extention := range imageExtentions {
		if strings.Contains(filename, extention) {
			return true
		}
	}
	return false
}

// opens given files, copies one into another (NOT EXPORTED !)
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
