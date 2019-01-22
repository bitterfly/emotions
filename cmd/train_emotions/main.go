package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bitterfly/emotions/emotions"
)

func readEmotion(filename string, k int) emotions.GaussianMixture {
	wf, _ := emotions.Read(filename, 0, 0.97)

	mfccs := emotions.MFCCs(wf, 13, 23)
	return emotions.GMM(mfccs, k)
}

func main() {
	k := 5

	if len(os.Args) < 3 {
		panic("go run main.go <outut_dir> <emotion1.wav [emotion2.wav...]>")
	}

	output_dir := os.Args[1]
	filenames := make([]string, len(os.Args)-2, len(os.Args)-2)
	egms := make([]emotions.EmotionGausianMixure, len(os.Args)-2, len(os.Args)-2)

	for i := 2; i < len(os.Args); i++ {

		filename := filepath.Base(os.Args[i])
		name := filename[0 : len(filename)-len(filepath.Ext(filename))]
		filenames[i-2] = filepath.Join(output_dir, name+".gmm")

		egms[i-2] = emotions.EmotionGausianMixure{
			Emotion: name,
			GM:      readEmotion(os.Args[i], k),
		}
	}

	for i := 0; i < len(filenames); i++ {
		bytes, err := json.Marshal(egms[i])
		if err != nil {
			panic(fmt.Sprintf("Error when marshaling %s\n", filenames[i]))
		}
		ioutil.WriteFile(filenames[i], bytes, 0644)
	}

	var bla emotions.EmotionGausianMixure
	m, _ := ioutil.ReadFile(filenames[0])
	err := json.Unmarshal(m, &bla)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", bla.Emotion)

}
