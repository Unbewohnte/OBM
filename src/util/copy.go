package util

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// opens given files, copies one into another
func CopyFile(srcPath, dstPath string) error {
	// open for reading (source)
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not open src file : %s", err))
	}
	defer srcFile.Close()

	// open for writing (destination) (create file, if it does not exist already)
	dstFile, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
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
