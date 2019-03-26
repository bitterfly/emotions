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

func testEmotion(emotion string, data [][]float64, egms []emotions.EmotionGausianMixure) {
	fmt.Printf("%s\t", emotion)

	emotions := make([]string, 0, len(egms))
	for _, egm := range egms {
		emotions = append(emotions, egm.Emotion)
	}

	counters := emotions.TestGMM(emotions, data, egms)

	keys := emotions.SortKeys(counters)
	for _, k := range keys {
		fmt.Printf("%d\t", counters[k])
	}
	fmt.Printf("\n")
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

	emotions := make([]string, 0, len(emotionFiles))
	for e := range emotionFiles {
		emotions = append(emotions, e)
	}

	sort.Strings(emotions)

	for _, emotion := range emotions {
		for _, file := range emotionFiles[emotion] {
			testEmotion(emotion, readEmotion(file), egms)
		}
	}
}
