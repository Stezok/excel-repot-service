package report

import (
	"io"
	"os"
)

func SafeCopyFile(to, from string) error {
	toFile, err := os.OpenFile(to, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer toFile.Close()

	fromFile, err := os.Open(from)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	_, err = io.Copy(toFile, fromFile)
	if err != nil {
		return err
	}

	return nil
}
