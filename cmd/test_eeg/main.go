package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/bitterfly/emotions/emotions"
)

func main() {
	if len(os.Args) < 5 {
		panic("go run main.go <classifier-type> <feature-type> <bucket-size> <gmm-dir> <input-file>\n<input-file>:<emotion>	<csv-file>")
	}

	//if bucket size is:
	// 0 take feature vectors every 200ms
	// 1 take only one vector for the whole file
	// n, n ≥ 2 take feature vector every n ms

	classifierType := os.Args[1]
	featureType := os.Args[2]
	bucketSize, err := strconv.Atoi(os.Args[3])
	if err != nil {
		panic(fmt.Sprintf("could not parse bucket-size argument: %s", os.Args[3]))
	}

	trainDir := os.Args[4]
	emotionFiles, _, _, err := emotions.ParseArgumentsFromFile(os.Args[5], false)

	frameLen := 200
	frameStep := 150

	if err != nil {
		panic(err)
	}
	switch classifierType {
	case "knn":
		err = emotions.ClassifyKNN(featureType, trainDir, bucketSize, frameLen, frameStep, emotionFiles)
		if err != nil {
			panic(err.Error())
		}
	case "gmm":
		err = emotions.ClassifyGMM(featureType, trainDir, bucketSize, frameLen, frameStep, emotionFiles)
	default:
		panic(fmt.Sprintf("Unknown classifier %s", classifierType))
	}

}
