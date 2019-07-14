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

// func testEmotion(emotion string, data [][]float64, egms []emotions.EmotionGausianMixure) (int, int, int) {
// 	fmt.Printf("%s\t", emotion)

// 	emotionNames := make([]string, 0, len(egms))
// 	for _, egm := range egms {
// 		emotionNames = append(emotionNames, egm.Emotion)
// 	}

// 	counters := emotions.TestGMM(emotionNames, data, egms)
// 	sum := 0
// 	keys := emotions.SortKeys(counters)
// 	for _, k := range keys {
// 		fmt.Printf("%d\t", counters[k])
// 		sum += counters[k]
// 	}
// 	fmt.Printf("\n")
// 	return correct(emotion, counters), counters[emotion], sum
// }

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
			boolCorrect, correctVector, sumVector := emotions.TestGMM(emotion, emotionTypes, readEmotion(file), egms)
			correctFiles[emotion] += boolCorrect
			correctVectors[emotion] += correctVector
			sumVectors[emotion] += sumVector
		}
	}

	fmt.Printf("\tCorrectFiles\tCorrectVectors\n")
	for _, emotion := range emotionTypes {
		fmt.Printf("%s\t%f\t%f\n", emotion, float64(correctFiles[emotion])/float64(len(emotionFiles[emotion])), float64(correctVectors[emotion])/float64(sumVectors[emotion]))
	}
}
