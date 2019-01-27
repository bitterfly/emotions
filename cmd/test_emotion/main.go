package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/bitterfly/emotions/emotions"
)

func readEmotion(filename string) [][]float64 {
	wf, _ := emotions.Read(filename, 0, 0.97)
	return emotions.MFCCs(wf, 13, 23)
}

func testEmotion(emotion string, coefficient [][]float64, egmms []emotions.EmotionGausianMixure) bool {
	k := len(egmms[0].GM)

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
	// fmt.Printf("======================\nEmotion: %s\n", emotion)
	for i, c := range counters {
		if c > max {
			max = c
			argmax = i
		}
		// fmt.Fprintf(os.Stderr, "%s: %d ", egmms[i].Emotion, c)
	}
	fmt.Fprintf(os.Stderr, "\nEmotion: %s Max: %s\n============================\n", emotion, egmms[argmax].Emotion)
	return strings.Contains(emotion, egmms[argmax].Emotion)
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
		panic("go run main.go <gmm-dir> <emotion_test_dir>")
	}

	gmmDir := os.Args[1]
	egms := getEGMs(gmmDir)

	name := os.Args[2]
	if filepath.Ext(name) == ".wav" {
		if testEmotion(name[0:len(name)-len(filepath.Ext(name))], readEmotion(name), egms) {
			fmt.Printf("100%%\n")
		} else {
			fmt.Printf("0%%\n")
		}
		os.Exit(0)
	}

	emotionDir := os.Args[2]

	files, err := ioutil.ReadDir(emotionDir)
	if err != nil {
		panic(err)
	}

	ctr := 0
	for _, file := range files {
		filename := file.Name()
		name := filename[0 : len(filename)-len(filepath.Ext(filename))]
		if testEmotion(name, readEmotion(path.Join(emotionDir, filename)), egms) {
			ctr++
		}
	}

	fmt.Printf("%f%%\n", float64(ctr)*100.0/float64(len(files)))
}
