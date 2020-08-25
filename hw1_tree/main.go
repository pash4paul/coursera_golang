package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	return printDir(out, file, printFiles, "")
}

func printDir(out io.Writer, file *os.File, printFiles bool, prefix string) error {
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		files, filesInfos, err := processDir(file, printFiles)
		if err != nil {
			return err
		}

		lastElementIndex := len(files) - 1
		nextPrefix := prefix + "│	"
		connector := "├───"

		for index, innerFile := range files {
			if index == lastElementIndex {
				connector = "└───"
				nextPrefix = prefix + "	"
			}

			outString := prefix + connector + filesInfos[innerFile].Name()
			if !filesInfos[innerFile].IsDir() {
				size := filesInfos[innerFile].Size()
				if size == 0 {
					outString += " (empty)"
				} else {
					outString = fmt.Sprintf("%s (%db)", outString, filesInfos[innerFile].Size())
				}
			}

			fmt.Fprintln(out, outString)

			err := printDir(out, innerFile, printFiles, nextPrefix)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func processDir(dir *os.File, printFiles bool) ([]*os.File, map[*os.File]os.FileInfo, error) {
	var files []*os.File
	filesInfos := make(map[*os.File]os.FileInfo)

	names, err := dir.Readdirnames(0)

	if err != nil {
		return files, filesInfos, err
	}

	sort.Strings(names)

	for _, name := range names {
		innerFile, err := os.Open(filepath.Join(dir.Name(), name))
		if err != nil {
			return files, filesInfos, err
		}

		fileInfo, err := innerFile.Stat()
		if err != nil {
			return files, filesInfos, err
		}

		if fileInfo.IsDir() || printFiles {
			files = append(files, innerFile)
			filesInfos[innerFile] = fileInfo
		}
	}

	return files, filesInfos, nil
}
