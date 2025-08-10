package filter

import (
	"bufio"
	"os"
	"slices"
)

var allWords map[string]struct{}

func InitWords(wordFile string) error {
	f, err := os.Open(wordFile)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		word := scanner.Text()
		if len(word) == 5 {
			allWords[word] = struct{}{}
		}
	}

	return nil
}

func match(word string, green [5]rune, yello map[rune][]int, greys map[rune]bool) bool {
	yelloMatch := make(map[rune]bool, len(yello))
	for i, c := range word {
		if _, grey := greys[c]; grey {
			return false
		}
		if green[i] != ' ' && c != green[i] {
			return false
		}
		if slices.Contains(yello[c], i) {
			return false
		}
		if _, ok := yello[c]; ok {
			yelloMatch[c] = true
		}
	}
	return len(yelloMatch) == len(yello)
}

func GetWords(green [5]rune, yello map[rune][]int, greys map[rune]bool) []string {
	words := []string{}
	for word := range allWords {
		if match(word, green, yello, greys) {
			words = append(words, word)
		}
	}
	return words
}

func init() {
	allWords = make(map[string]struct{})
}
