package main

import (
	"os"

	"github.com/bitterfly/emotions/emotions"
)

func main() {
	if len(os.Args) < 3 {
		panic("go run main.go <speech-gmm-dir> <eeg-gmm-dir> <input-file>\n<input-file>:<emotion>	<csv-file>")
	}

	speechDir := os.Args[1]
	eegDir := os.Args[2]
	speechFiles, eegFiles, err := emotions.ParseArgumentsFromFile(os.Args[3], true)
	if err != nil {
		panic(err)
	}
	emotions.ClassifyGMMBothConcat(speechDir, speechFiles, eegDir, eegFiles)
}
