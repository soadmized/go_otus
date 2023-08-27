package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	currFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0o444)
	if err != nil {
		return err
	}

	defer currFile.Close()

	err = validateSizeAndOffset(currFile, offset)
	if err != nil {
		return err
	}

	newFile, err := os.Create(toPath)
	if err != nil {
		return err
	}

	defer newFile.Close()

	limit, err = prepareLimit(currFile, limit)
	if err != nil {
		return err
	}

	reader := io.NewSectionReader(currFile, offset, limit)

	err = copyWithProgressBar(newFile, reader)
	if err != nil {
		return err
	}

	return nil
}

func validateSizeAndOffset(file *os.File, offset int64) error {
	fi, err := file.Stat()
	if err != nil {
		return err
	}

	size := fi.Size()
	if offset > size {
		return ErrOffsetExceedsFileSize
	}
	if size == 0 {
		return ErrUnsupportedFile
	}

	return nil
}

func prepareLimit(currFile *os.File, limit int64) (int64, error) {
	fi, err := currFile.Stat()
	if err != nil {
		return 0, err
	}

	size := fi.Size()
	if (limit > size) || (limit == 0) {
		return size, nil
	}

	return limit, nil
}

func copyWithProgressBar(writer io.Writer, reader *io.SectionReader) error {
	bar := pb.Full.Start64(reader.Size())
	barReader := bar.NewProxyReader(reader)

	_, err := io.Copy(writer, barReader)
	if err != nil {
		return err
	}

	bar.Finish()

	return nil
}
