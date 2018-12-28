package main

import (
	"fmt"
	"os"

	"github.com/bitterfly/emotions/fourier"
)

func main() {
	// dirname := os.Args[1]
	// files, _ := ioutil.ReadDir(dirname)

	// indices := make([]int, len(files), len(files))
	// mfccs := make([][]float64, 0, len(files)*1000)
	// names := make([]string, len(files), len(files))
	// for i, f := range files {
	// 	names[i] = f.Name()[0 : len(f.Name())-4]
	// 	wf, err := fourier.Read(path.Join(dirname, f.Name()), 0, 0.97)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	mfccs = append(mfccs, fourier.MFCCs(wf, 13, 23)...)
	// 	indices[i] = len(mfccs) - 1
	// }

	// fmt.Printf("%d\n", len(mfccs))
	// fmt.Printf("Kmeans: \n")
	points := [][]float64{
		[]float64{151700, 351102},
		[]float64{155799, 354358},
		[]float64{142857, 352716},
		[]float64{152726, 349144},
		[]float64{151008, 349692},
	}

	c, ms := fourier.Kmeans(points, 1)

	pf, _ := os.Create("/tmp/points.csv")
	cf, _ := os.Create("/tmp/centroids.csv")
	defer pf.Close()
	defer cf.Close()

	fmt.Fprintf(pf, "X, Y, Z\n")
	fmt.Fprintf(cf, "X, Y, Z\n")

	for i, ms := range ms {
		fmt.Fprintf(pf, "%f, %f, %d\n", points[i][0], points[i][1], ms.GetCluster())
	}

	for i, cc := range c {
		fmt.Fprintf(cf, "%f, %f, %d\n", cc[0], cc[1], i)
	}
}
