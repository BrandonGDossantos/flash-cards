package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type allJSON struct {
	all []fileJSON
}

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
	// Exit /text directory to main directory
	fmt.Println(fileLines[0])
	os.Chdir("../")
	return fileLines, nil
}

// Q|A|Q|A|Q|A|Q|A
// 0|1|2|3|4|5|6|7
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
	// prettyMap(cards)
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
