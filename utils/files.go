package utils

import (
	"fmt"
	"io"
	"os"
)

func Copy(src, dst string) error {
	sourceFilesStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFilesStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not correct", src)
	}

	source, err := os.Open(src)

	if err != nil {
		return err
	}

	defer source.Close()

	destination, err := os.Create(dst)

	if err != nil {
		return err
	}

	defer destination.Close()
	_, err = io.Copy(destination, source)

	return err
}
