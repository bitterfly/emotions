package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/bitterfly/emotions/emotions"
)

func main() {
	if len(os.Args) < 5 {
		panic("go run main.go <bucket-size> <speech-gmm-dir> <eeg-gmm-dir> <input-file>\n<input-file>:<emotion>	<csv-file>")
	}

	//if bucket size is:
	// 0 take feature vectors every 200ms
	// 1 take only one vector for the whole file
	// n, n ≥ 2 take feature vector every n ms

	bucketSize, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(fmt.Sprintf("could not parse bucket-size argument: %s", os.Args[3]))
	}

	speechDir := os.Args[2]
	eegDir := os.Args[3]
	speechFiles, eegFiles, err := emotions.ParseArgumentsFromFile(os.Args[4], true)

	frameLen := 200
	frameStep := 150

	emotions.ClassifyGMMBoth(bucketSize, frameLen, frameStep, speechDir, speechFiles, eegDir, eegFiles)
}
