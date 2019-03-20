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

	bucketSize, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(fmt.Sprintf("could not parse bucket-size argument: %s", os.Args[1]))
	}

	trainFile := os.Args[2]
	emotionFiles, _, err := emotions.ParseArgumentsFromFile(os.Args[3], false)

	if err != nil {
		panic(err)
	}

	err = emotions.KNN(bucketSize, 200, 150, trainFile, emotionFiles)
	if err != nil {
		panic(err.Error())
	}

}
