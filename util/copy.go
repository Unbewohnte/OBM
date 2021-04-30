package util

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// opens given files, copies one into another
func CopyFile(src, dst string) error {
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
