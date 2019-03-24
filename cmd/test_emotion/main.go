package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"

	"github.com/bitterfly/emotions/emotions"
)

func readEmotion(filename string) [][]float64 {
	wf, _ := emotions.Read(filename, 0.01, 0.97)
	return emotions.MFCCs(wf, 13, 23)
}

func testEmotion(emotion string, coefficient [][]float64, egmms []emotions.EmotionGausianMixure) {
	k := len(egmms[0].GM)

	counters := make(map[string]int)
	for _, m := range coefficient {
		max := math.Inf(-42)
		argmax := -1
		for i, egmm := range egmms {
			currEmotion := emotions.EvaluateVector(m, k, egmm.GM)
			if currEmotion > max {
				max = currEmotion
				argmax = i
			}
		}
		counters[egmms[argmax].Emotion]++
	}

	fmt.Printf("%s\t", emotion)

	keys := emotions.SortKeys(counters)

	for _, k := range keys {
		fmt.Printf("%d\t", counters[k])
	}
	fmt.Printf("\n")

}

func getEGMs(dirname string) []emotions.EmotionGausianMixure {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Fatal(err)
	}

	egms := make([]emotions.EmotionGausianMixure, len(files), len(files))
	for i, f := range files {
		bytes, _ := ioutil.ReadFile(filepath.Join(dirname, f.Name()))
		err := json.Unmarshal(bytes, &egms[i])
		if err != nil {
			panic(err)
		}
	}

	return egms
}

func main() {
	if len(os.Args) < 3 {
		panic("go run main.go <gmm-dir> <input-file>\n<input-file>:<emotion> <wav-file>")
	}

	gmmDir := os.Args[1]
	egms := getEGMs(gmmDir)

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
