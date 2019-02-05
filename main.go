package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// fileJSON is the structure for each text to JSON conversion
type fileJSON struct {
	Title string
	Cards map[string]string
}

// convertFile opens a file, reads each line, and returns
// a list of every line
func convertFile(f string) ([]string, error) {
	var lines []string
	openFile, err := os.Open(f)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(openFile)
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	openFile.Close()
	return lines, nil
}

// readDir traverses through every file in a directory and calls convertFile() on each
// returning a list of lists that contain each line in every file in the directory.
func readDir(dir string) ([][]string, error) {
	f, err := os.Open(dir)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	files, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}
	// Enter the /text directory to access all .txt files
	os.Chdir(dir)
	var fileLines [][]string
	for _, file := range files {
		fLines, err := convertFile(file.Name())
		if err != nil {
			return nil, err
		}
		fileLines = append(fileLines, fLines)
	}
	// Exit /text directory
	os.Chdir("../")
	return fileLines, nil
}

// writeJSON reads the 2D slice of file text and writes it out to the JSON structs
func writeJSON(readLines [][]string) {
	var a []fileJSON
	for _, readLine := range readLines {
		var cards = make(map[string]string)
		for i := 1; i < len(readLine); i += 2 {
			cards[readLine[i]] = readLine[i+1]
		}
		f := fileJSON{
			Title: readLine[0],
			Cards: cards,
		}
		a = append(a, f)
	}
	var jsonData []byte
	jsonData, err := json.MarshalIndent(a, "", "    ")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(jsonData))
}

func main() {
	readLines, err := readDir("./text")
	if err != nil {
		log.Fatal(err)
	}
	writeJSON(readLines)
}
