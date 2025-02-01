package main

import (
	"errors"
	"github.com/cheggaaa/pb/v3"
	"io"
	"log"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrOffsetNegativeOffset  = errors.New("negative offset is not acceptable")
	ErrOffsetNegativeLimit   = errors.New("negative limit is not acceptable")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Place your code here.
	fileFrom, err := os.Open(fromPath)

	if err != nil {
		return err
	}

	defer func() {
		if err := fileFrom.Close(); err != nil {
			log.Printf("Failed to close file %s: %v\n", fromPath, err)
		}
	}()

	fileFromInfo, err := os.Stat(fromPath)
	if err != nil {
		return err
	}

	if offset > fileFromInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	if offset < 0 {
		return ErrOffsetNegativeOffset
	}

	if limit < 0 {
		return ErrOffsetNegativeLimit
	}

	if fileFromInfo.Mode()&os.ModeDevice != 0 {
		return ErrUnsupportedFile
	}

	fileTo, err := os.Create(toPath)

	if err != nil {
		return err
	}

	defer func() {
		if err := fileTo.Close(); err != nil {
			log.Printf("Failed to close file %s: %v\n", to, err)
		}
	}()

	fileSize := fileFromInfo.Size() - offset

	var bytesToCopy int64
	if limit == 0 || limit > fileSize {
		bytesToCopy = fileSize
	} else {
		bytesToCopy = limit
	}

	bar := pb.StartNew(int(bytesToCopy))

	_, err = fileFrom.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	var moved int64
	for moved < fileSize {
		written, err := io.CopyN(fileTo, fileFrom, bytesToCopy)
		moved += written

		bar.Add(int(moved))

		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if moved >= limit && limit != 0 {
			break
		}
	}

	bar.Finish()

	return nil
}
