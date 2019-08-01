package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/bitterfly/emotions/emotions"
)

func readEmotion(filename string) [][]float64 {
	wf, _ := emotions.Read(filename, 0.01, 0.97)
	return emotions.MFCCs(wf, 13, 23)
}

func correct(emotion string, counters map[string]int) int {
	maxV := 0
	maxE := ""
	for e, v := range counters {
		if v > maxV {
			maxV = v
			maxE = e
		}
	}
	if maxE == emotion {
		return 1
	}
	return 0
}

func main() {
	if len(os.Args) < 3 {
		panic("go run main.go <gmm-dir> <input-file>\n<input-file>:<emotion> <wav-file>")
	}

	gmmDir := os.Args[1]
	egms, err := emotions.GetEGMs(gmmDir)
	if err != nil {
		panic(err)
	}

	emotionFiles, _, err := emotions.ParseArgumentsFromFile(os.Args[2], false)
	if err != nil {
		panic(err)
	}

	emotionTypes := make([]string, 0, len(emotionFiles))
	for e := range emotionFiles {
		emotionTypes = append(emotionTypes, e)
	}

	sort.Strings(emotionTypes)
	correctFiles := make(map[string]int, len(emotionTypes))
	correctVectors := make(map[string]int, len(emotionTypes))
	sumVectors := make(map[string]int, len(emotionTypes))

	for _, emotion := range emotionTypes {
		for _, file := range emotionFiles[emotion] {
			boolCorrect, vectors, sumVector := emotions.TestGMM(emotion, emotionTypes, readEmotion(file), egms)
			correctFiles[emotion] += boolCorrect
			correctVectors[emotion] += vectors[emotion]
			sumVectors[emotion] += sumVector
		}
	}

	fmt.Printf("\tCorrectFiles\tCorrectVectors\n")
	for _, emotion := range emotionTypes {
		fmt.Printf("%s\t%f\t%f\n", emotion, float64(correctFiles[emotion])/float64(len(emotionFiles[emotion])), float64(correctVectors[emotion])/float64(sumVectors[emotion]))
	}
}
