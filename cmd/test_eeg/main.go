package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/bitterfly/emotions/emotions"
)

func main() {
	if len(os.Args) < 3 {
		panic("go run main.go <bucket-size> <eeg-train-file> <input-file>\n<input-file>:<emotion>	<csv-file>")
	}

	//if bucket size is:
	// 0 take feature vectors every 200ms
	// 1 take only one vector for the whole file
	// n, n ≥ 2 take feature vector every n ms

	classifierType := os.Args[1]
	bucketSize, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(fmt.Sprintf("could not parse bucket-size argument: %s", os.Args[1]))
	}

	trainDir := os.Args[3]
	emotionFiles, _, err := emotions.ParseArgumentsFromFile(os.Args[4], false)

	frameLen := 200
	frameStep := 150

	if err != nil {
		panic(err)
	}
	switch classifierType {
	case "knn":
		err = emotions.ClassifyKNN(trainDir, bucketSize, frameLen, frameStep, emotionFiles)
		if err != nil {
			panic(err.Error())
		}
	case "gmm":
		err = emotions.ClassifyGMM(trainDir, bucketSize, frameLen, frameStep, emotionFiles)
	default:
		panic(fmt.Sprintf("Unknown classifier %s", classifierType))
	}

}
