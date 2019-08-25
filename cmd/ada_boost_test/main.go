package main

import (
	"os"

	"github.com/bitterfly/emotions/emotions"
)

func main() {
	if len(os.Args) < 4 {
		panic("go run main.go <speech-gmm-dir> <eeg-gmm-dir> <input-file>\n<input-file>:<emotion>	<csv-file>")
	}

	speechDir := os.Args[1]
	eegDir := os.Args[2]
	speechFiles, eegFiles, _, err := emotions.ParseArgumentsFromFile(os.Args[3], true)

	if err != nil {
		panic(err)
	}

	bucketSize := 0
	frameLen := 200
	frameStep := 150

	emotions.ClassifyGMMBoth(bucketSize, frameLen, frameStep, speechDir, speechFiles, eegDir, eegFiles)
}
