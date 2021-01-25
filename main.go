package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

func main() {
	var firstfilePath string
	var secondfilePath string

	// input two path
	fmt.Print("input first file path: ")
	fmt.Scan(&firstfilePath)
	fmt.Print("input second file path: ")
	fmt.Scan(&secondfilePath)

	// openned files
	firstFile, firstFileErr := os.Open(firstfilePath)
	secondFile, secondFileErr := os.Open(secondfilePath)

	defer firstFile.Close()
	defer secondFile.Close()

	// check two file
	status, errs := checkStatus(firstFileErr, secondFileErr)

	if status {
		// read two json file
		var firstDecoder *json.Decoder = json.NewDecoder(firstFile)
		var secondDecoder *json.Decoder = json.NewDecoder(secondFile)

		var firstJsonData map[string]interface{}
		var secondJsonData map[string]interface{}
		var mapSlice map[string][]string = map[string][]string{}

		firstDecoderErr := firstDecoder.Decode(&firstJsonData)
		secondDecoderErr := secondDecoder.Decode(&secondJsonData)

		// check two json format
		status, errs = checkStatus(firstDecoderErr, secondDecoderErr)

		if status {
			for key, val := range firstJsonData {
				if secondJsonData[key] != nil {

					if secondJsonData[key] != val {
						mapSlice["change"] = append(mapSlice["change"], key)
					} else {
						mapSlice["same"] = append(mapSlice["same"], key)
					}

					delete(secondJsonData, key)
				} else {
					mapSlice["remove"] = append(mapSlice["remove"], key)
				}
			}

			for key, _ := range secondJsonData {
				mapSlice["add"] = append(mapSlice["add"], key)
			}

			writeResult(mapSlice)
			fmt.Println("\nWork completed!")
		} else {
			fmt.Println("\n!!!Error!!! Open error.txt file")
			writerFile("error.txt", errs)
		}
	} else {
		fmt.Println("\n!!!Error!!! Open error.txt file")
		writerFile("error.txt", errs)
	}
}

func checkStatus(errors ...error) (status bool, ans []string) {
	status = true

	for _, val := range errors {
		if val != nil {
			ans = append(ans, val.Error())
			status = false
		}
	}

	return
}

func writerFile(path string, errorTexts []string) {
	file, err := os.Create(path)

	writer := bufio.NewWriter(file)

	defer file.Close()

	if err == nil {
		for _, val := range errorTexts {
			writer.WriteString(val + "\n")
		}
	}

	writer.Flush()
}

func writeResult(mapSlice map[string][]string) {
	for key, val := range mapSlice {
		sort.Strings(val)
		path := key + ".txt"
		writerFile(path, val)
	}
}
