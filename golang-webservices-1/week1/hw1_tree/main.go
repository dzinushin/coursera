package main

import (
	"errors"
	"fmt"
	"io"
	"os"
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
	//fmt.Println("> dirTree path: ", path)
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		//fmt.Println("closing file")
		_ = file.Close()
	}(file)

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	isDir := fileInfo.IsDir()
	//fmt.Printf("file name: %s, isDir: %v\n", file.Name(), isDir)
	if !isDir {
		return errors.New("only dirs allowed")
	}
	return listDir(out, path, printFiles, "")
}

func listDir(out io.Writer, path string, printFiles bool, prefix string) error {

	dirEntries, err := getDirEntries(path, printFiles)
	if err != nil {
		return err
	}

	//fmt.Printf("path: %v dirEntries len: %v\n", path, len(dirEntries))
	totalEntries := len(dirEntries)
	for idx, dirEntry := range dirEntries {
		isLastItem := idx == totalEntries-1
		printDirEntry(out, dirEntry, prefix, isLastItem)

		isDir := dirEntry.IsDir()
		if isDir {
			subDirPrefix := prefix
			if isLastItem {
				subDirPrefix += "\t"
			} else {
				subDirPrefix += "│\t"
			}
			err := listDir(out, path+"/"+dirEntry.Name(), printFiles, subDirPrefix)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getDirEntries(path string, printFiles bool) ([]os.DirEntry, error) {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return []os.DirEntry{}, err
	}
	sort.Slice(dirEntries, func(i, j int) bool { return dirEntries[i].Name() < dirEntries[j].Name() })

	if printFiles {
		return dirEntries, nil
	}
	var onlyDirEntries = make([]os.DirEntry, 0)
	for _, entry := range dirEntries {
		if entry.IsDir() {
			onlyDirEntries = append(onlyDirEntries, entry)
		}
	}
	return onlyDirEntries, nil
}

func printDirEntry(out io.Writer, entry os.DirEntry, dirPrefix string, isLastDirEntry bool) {
	var itemPrefix string
	if isLastDirEntry {
		itemPrefix = "└───"
	} else {
		itemPrefix = "├───"
	}
	isDir := entry.IsDir()

	var itemStr string
	if isDir {
		itemStr = entry.Name()
	} else {
		info, _ := entry.Info()
		size := info.Size()
		var itemStrSuffix string
		if size == 0 {
			itemStrSuffix = "(empty)"
		} else {
			itemStrSuffix = fmt.Sprintf("(%db)", size)
		}
		itemStr = fmt.Sprintf("%s %s", entry.Name(), itemStrSuffix)
	}

	fmt.Fprintf(out, "%s%s%s\n", dirPrefix, itemPrefix, itemStr)
}
