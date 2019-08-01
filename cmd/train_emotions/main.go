package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"

	"github.com/bitterfly/emotions/emotions"
)

func main() {
	if len(os.Args) < 3 {

		panic("go run main.go <k> <dir-template> <input_file>\n<input_file>: <emotion>	<wav_file>")
	}

	outputDirIndex := 2
	maxK := -1

	k, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic("Must provide k")
	}
	if _, err := strconv.Atoi(os.Args[2]); err != nil {
		outputDirIndex = 2
		maxK = k
	} else {
		maxK, _ = strconv.Atoi(os.Args[2])
		outputDirIndex = 3
	}

	outputDir := os.Args[outputDirIndex]
	emotionFiles, _, err := emotions.ParseArgumentsFromFile(os.Args[outputDirIndex+1], false)

	if err != nil {
		panic(err)
	}

	for j := k; j <= maxK; j++ {
		if _, err := os.Stat(fmt.Sprintf("%s_k%d", outputDir, j)); os.IsNotExist(err) {
			os.Mkdir(fmt.Sprintf("%s_k%d", outputDir, j), 0775)
		}
	}

	emotionTypes := make([]string, 0, len(emotionFiles))
	for e := range emotionFiles {
		emotionTypes = append(emotionTypes, e)
	}

	sort.Strings(emotionTypes)

	for _, emotion := range emotionTypes {
		files := emotionFiles[emotion]
		mfccs := emotions.ReadSpeechFeatures(files)
		for j := k; j <= maxK; j++ {
			fmt.Fprintf(os.Stderr, "%s %d\n", emotion, j)
			fmt.Printf("k: %d SpeechFiles: %v\n", j, files)

			egm := emotions.EmotionGausianMixure{
				Emotion: emotion,
				GM:      emotions.GMM(mfccs, j),
			}

			bytes, err := json.Marshal(egm)
			if err != nil {
				panic(fmt.Sprintf("Error when marshaling %s with k %d\n", emotion, j))
			}

			filename := path.Join(fmt.Sprintf("%s_k%d", outputDir, j), fmt.Sprintf("%s.gmm", emotion))
			fmt.Fprintf(os.Stderr, filename)
			ioutil.WriteFile(filename, bytes, 0644)
		}
	}
}
