package util

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

// creates a complete black image file
func CreateBlackBG(width, height uint) error {
	bgfile, err := os.Create("blackBG.png")
	if err != nil {
		return errors.New(fmt.Sprintf("Could not create black background file : %s", err))
	}
	image := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
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
