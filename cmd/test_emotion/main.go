package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"

	"github.com/bitterfly/emotions/emotions"
)

func readEmotion(filename string) [][]float64 {
	wf, _ := emotions.Read(filename, 0, 0.97)
	return emotions.MFCCs(wf, 13, 23)
}

func testEmotion(k int, emotion string, coefficient [][]float64, egmms []emotions.EmotionGausianMixure) {
	counters := make([]int, len(egmms), len(egmms))
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
		counters[argmax]++
	}

	max := -1
	argmax := -1
	fmt.Printf("======================\nEmotion: %s\n", emotion)
	for i, c := range counters {
		if c > max {
			max = c
			argmax = i
		}
		fmt.Printf("%s: %d ", egmms[i].Emotion, c)
	}
	fmt.Printf("\nMax: %s\n============================\n", egmms[argmax].Emotion)
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
	k := 5

	if len(os.Args) < 3 {
		panic("go run main.go <gmm-dir> <emotion1.wav [emotion2.wav...]>")
	}

	gmmDir := os.Args[1]
	egms := getEGMs(gmmDir)

	for i := 2; i < len(os.Args); i++ {
		filename := filepath.Base(os.Args[i])
		name := filename[0 : len(filename)-len(filepath.Ext(filename))]
		testEmotion(k, name, readEmotion(os.Args[i]), egms)
	}
}
