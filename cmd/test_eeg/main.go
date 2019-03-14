package main

import (
	"os"

	"github.com/bitterfly/emotions/emotions"
)

func main() {
	if len(os.Args) < 6 {
		panic("go run main.go <eeg-train-file> --eeg-positive eeg_pos1.txt [eeg_pos2.txt...] --eeg-negative eeg_neg1.txt [eeg_neg2.txt...] --eeg-neutral eeg_neu1.txt [eeg_neu2.txt...]")
	}

	trainFile := os.Args[1]
	emotionFiles := emotions.ParseArguments(os.Args[2:])

	err := emotions.KNN(trainFile, emotionFiles)
	if err != nil {
		panic(err.Error())
	}

}
