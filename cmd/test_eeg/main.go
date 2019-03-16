package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/bitterfly/emotions/emotions"
)

func main() {
	if len(os.Args) < 6 {
		panic("go run main.go <bucket-size> <eeg-train-file> --eeg-positive eeg_pos1.txt [eeg_pos2.txt...] --eeg-negative eeg_neg1.txt [eeg_neg2.txt...] --eeg-neutral eeg_neu1.txt [eeg_neu2.txt...]")
	}

	//if bucket size is:
	// 0 take feature vectors every 200ms
	// 1 take only one vector for the whole file
	// n, n ≥ 2 take feature vector every n ms

	bucketSize, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(fmt.Sprintf("could not parse bucket-size argument: %s", os.Args[1]))
	}

	trainFile := os.Args[2]
	emotionFiles := emotions.ParseArguments(os.Args[3:])

	err = emotions.KNN(bucketSize, 200, 150, trainFile, emotionFiles)
	if err != nil {
		panic(err.Error())
	}

}
