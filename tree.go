package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
)

const CHUNK = 8192

func writeToFile(f *os.File, filename string, md5Value string) {
	// Get file info
	info, err := os.Stat(filename)
	if err != nil {
		err = errors.New(fmt.Sprintf("Cannot access file:", filename))
		return
	}

	// Get the file size
	fileSize := info.Size()

	_, err = f.WriteString(fmt.Sprintf("%s,  %s, %d\n", filename, md5Value, fileSize))
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
	fmt.Println("Usage: go run tree.go <path> [ignore_patterns]")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	directory := os.Args[1]
	ignorePatterns := make([]string, 0)
	for i := 2; i < len(os.Args); i++ {
		ignorePatterns = append(ignorePatterns, os.Args[i])
	}

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

		matched := false
		for _, pattern := range ignorePatterns {
			m, err := filepath.Match(path, pattern)
			if err != nil {
				fmt.Printf("filepath.Match returned %v\n", err)
			} else {
				if m == true {
					matched = true
				}
			}
		}

		// ignore
		if matched == true {
			if f.IsDir() {
				return filepath.SkipDir
			} else {
				return nil
			} 
		}
		
		if f.IsDir() {
			writeToFile(outputFile, path, "")
		} else {
			md5Value, err := getFileMd5(path)
			if err != nil {
				fmt.Printf("getFileMd5 returned %v\n", err)
			} else {
				writeToFile(outputFile, path, md5Value)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
}
