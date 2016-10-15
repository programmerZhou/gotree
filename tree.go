package main

import (
	"fmt"
	"math"
	"os"
	"errors"
	"path/filepath"
	"crypto/md5"
	"io"
)

const CHUNK = 8192

func writeToFile(f *os.File, filename string, md5Value string) {
	_, err := f.WriteString(fmt.Sprintf("%s    %s\n", filename, md5Value))
	if err != nil {
		fmt.Printf("Write file fail: %v", err)
	}
}

func getFileMd5(filename string) (md5Value string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		err = errors.New(fmt.Sprintf("Cannot find file: %s", filename))
		return
	}

	defer file.Close()

	// Get file info
	info, err := file.Stat()
	if err != nil {
		err = errors.New(fmt.Sprintf("Cannot access file:", filename))
		return
	}

	// Get the file size
	fileSize := info.Size()

	// Calculate the number of blocks
	blocks := uint64(math.Ceil(float64(fileSize) / float64(CHUNK)))

	// Start hash
	hash := md5.New()

	// Check each block
	for i := uint64(0); i < blocks; i++ {
		// Calculate block size
		blockSize := int(math.Min(CHUNK, float64(fileSize-int64(i*CHUNK))))

		// Make a buffer
		buf := make([]byte, blockSize)

		// Read to the buffer
		file.Read(buf)

		// Write from the buffer
		io.WriteString(hash, string(buf))
	}

	// Get md5Sum value
	md5Value = fmt.Sprintf("%x", hash.Sum(nil))

	return
}

func usage() {
	fmt.Println("go run tree.go path")
}

func main() {
	if len(os.Args) != 2 {
		usage()
		return
	}

	directory := os.Args[1]

	fileInfo, err := os.Stat(directory)
	if err != nil {
		fmt.Printf("Cannot accect this path: %s\n", directory)
		return
	} else {
		if !fileInfo.IsDir() {
			fmt.Printf("This path is not directory\n")
			return
		}
	}

	// Open output file
	outputFile, err := os.OpenFile("tree.log", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Open output file fail", err)
		return
	}
	defer outputFile.Close()

	// Traversal the directory
	err = filepath.Walk(directory, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return nil
		}
		if f.IsDir() {
			return nil
		}

		md5Value, err := getFileMd5(path)
		if err != nil {
			fmt.Printf("getFileMd5 returned %v\n", err)
		} else {
			writeToFile(outputFile, path, md5Value)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
}
