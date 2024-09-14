package main

import (
	"bufio"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func readFile(fileName string) string {
	data, err := os.ReadFile(fileName)
	check(err)

	return string(data)
}

func writeToFile(fileName string, data string) {
	f, err := os.Create(fileName)
	check(err)
	defer f.Close()

	w := bufio.NewWriter(f)
	w.WriteString(data)

	w.Flush()
}

func containsMultiple(text string, substrings []string) bool {
	for _, sub := range substrings {
		if strings.Contains(text, sub) {
			return true
		}
	}
	return false
}
