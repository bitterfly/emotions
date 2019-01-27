package emotions

import (
	"fmt"
	"math"
	"os"
)

const FRAME_IN_MS = 25
const STEP_IN_MS = 10

func rectangular(x int, n int) float64 {
	return 1
}

func hanning(x int, n int) float64 {
	return 0.5 - 0.5*math.Cos(float64(x*2)*math.Pi/float64(n-1))
}

// so when x == n-1 -> 2pi
func hamming(x int, n int) float64 {
	return 0.54 - 0.46*math.Cos(float64(x*2)*math.Pi/float64(n-1))
}

func window(content []float64, windowFunction func(int, int) float64) {
	for x, y := range content {
		content[x] = y * windowFunction(x, len(content))
	}
}

func hammingWindow(content []float64) {
	for x, y := range content {
		content[x] = y * hamming(x, len(content))
	}
}

func hanningWindow(content []float64) {
	for x, y := range content {
		content[x] = y * hanning(x, len(content))
	}
}

func rectangularWindow(content []float64) {
	for x, y := range content {
		content[x] = y
	}
}

// CutSliceIntoFrames receives float array and cuts it into frames
// It takes approximately 25mS of real data and then pads it with zeroes to the first power of two
// It takes frames with len 25ms with step of 10ms and applies a window function to the frame (hamming)
func CutSliceIntoFrames(data []float64, sampleRate uint32) [][]float64 {
	realSamplesPerFrame := int((FRAME_IN_MS / 1000.0) * float64(sampleRate))

	samplesPerFrame := FindClosestPower(int(realSamplesPerFrame))
	step := int((STEP_IN_MS / 1000.0) * float64(sampleRate))

	fmt.Fprintf(os.Stderr, "Slice len: %d\n", len(data))
	fmt.Fprintf(os.Stderr, "Real samples per frame for 25ms: %d\n", realSamplesPerFrame)
	fmt.Fprintf(os.Stderr, "Samples per frame: %d\nStep: %d\n", samplesPerFrame, step)

	fmt.Fprintf(os.Stderr, "Which is %f ms long\n", 1000.0*float64(samplesPerFrame)/float64(sampleRate))

	numFrames := (len(data) - realSamplesPerFrame) / step
	fmt.Fprintf(os.Stderr, "Num frames: %d\n", numFrames)
	frames := make([][]float64, numFrames, numFrames)

	frame := 0
	for i := 0; i < len(data)-realSamplesPerFrame-step; i += step {
		frames[frame] = sliceCopyWithWindow(data, i, i+realSamplesPerFrame, samplesPerFrame)

		frame++
	}
	return frames
}

func sliceCopyWithWindow(first []float64, from, to, length int) []float64 {
	second := make([]float64, length, length)
	copy(second, first[from:Min(to, len(first))])
	hammingWindow(second[0 : Min(to, len(first))-from])
	return second
}

//CutWavFileIntoFrames takes a wavfiles and cuts it into frames
func CutWavFileIntoFrames(wf WavFile) [][]float64 {
	return CutSliceIntoFrames(wf.data, wf.sampleRate)
}
