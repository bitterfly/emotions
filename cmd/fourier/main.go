package main

import "github.com/bitterfly/emotions/fourier"

func main() {
	// filename := os.Args[1]
	// wf, _ := fourier.Read(filename, 0, 0.97)

	// mfccs := fourier.MFCCs(wf, 13, 23)
	// doubles := fourier.MFCCcDouble(mfccs)

	// for i, d := range doubles {
	// 	fmt.Printf("%d %v\n", i, d)
	// 	fmt.Printf("\n")
	// }

	points := [][]float64{
		[]float64{1, 1},
		[]float64{1.5, 2},
		[]float64{3, 4},
		[]float64{5, 7},
		[]float64{3.5, 5},
		[]float64{4.5, 5},
		[]float64{3.5, 4.5},
	}

	fourier.Kmeans(points, 2)
}
