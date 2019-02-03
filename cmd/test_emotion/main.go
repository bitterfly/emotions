package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitterfly/emotions/emotions"
)

func readEmotion(emotion string, filename string) [][]float64 {
	wf, _ := emotions.Read(filename, 0.01, 0.97)
	mfccs := emotions.MFCCs(wf, 13, 23)
	for i := range mfccs {
		mfccs[i] = append(mfccs[i], emotions.GetValence(emotion, 0))
	}
	return mfccs
}

func testEmotion(emotion string, coefficient [][]float64, egmms []emotions.EmotionGausianMixure) bool {
	k := len(egmms[0].GM)

	counters := make([]int, len(egmms), len(egmms))
	for _, m := range coefficient {
		max := math.Inf(-42)
		argmax := -1
		for i, egmm := range egmms {
			fmt.Printf("%v\n", m)
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
	for i, c := range counters {
		if c > max {
			max = c
			argmax = i
		}
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
		panic("go run main.go <gmm-dir> <-emotion1> <emotion1.wav2 emotion1.wav2...> [<emotion2> <emotion2.wav1...] ")
	}

	gmmDir := os.Args[1]
	egms := getEGMs(gmmDir)

	emotionFiles := emotions.ParseArguments(os.Args[2:])

	allctr := 0
	allfiles := 0
	for emotion, files := range emotionFiles {
		ctr := 0
		allfiles += len(files)
		for _, file := range files {
			if testEmotion(emotion, readEmotion(emotion, file), egms) {
				ctr++
				allctr++
			}
		}
		fmt.Printf("%s %f%%\n", emotion, float64(ctr)*100.0/float64(len(files)))
	}
	fmt.Printf("All: %f%%\n", float64(allctr)*100.0/float64(allfiles))
}
