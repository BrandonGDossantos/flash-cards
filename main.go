package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type fileJSON struct {
	Title string
	Cards map[string]string
}

func prettyMap(m map[string]string) {
	// Creat tmp file to store the maps print out
	file, err := os.Create("tmp")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()
	for k, v := range m {
		fmt.Fprintf(file, "Question: %s\nAnswer: %s\n", k, v)
	}
}

func readFiles(dir string) ([]string, error) {
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
	var lines []string
	for _, file := range files {
		openFile, err := os.Open(file.Name())
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
	}
	// Exit /text directory to main directory
	os.Chdir("../")
	return lines, nil
}

// Q|A|Q|A|Q|A|Q|A
// 0|1|2|3|4|5|6|7
func writeJSON(readLines []string) {
	var cards = make(map[string]string)
	for i := 1; i < len(readLines); i += 2 {
		cards[readLines[i]] = readLines[i+1]
	}
	f := fileJSON{
		Title: readLines[0],
		Cards: cards,
	}
	// prettyMap(cards)
	var jsonData []byte
	jsonData, err := json.MarshalIndent(f, "", "    ")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(jsonData))
}

func main() {
	readLines, err := readFiles("./text")
	if err != nil {
		log.Fatal(err)
	}

	writeJSON(readLines)
}
