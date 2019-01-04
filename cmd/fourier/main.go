package main

import (
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
	// points := [][]float64{
	// 	[]float64{1, 0, 1.002},
	// 	[]float64{0, 1, 1.02},
	// 	[]float64{5.23, 5.22, 5.10},
	// 	[]float64{6.23, 6.4, 6.33},
	// 	[]float64{7.11, 7.22, 7.13},
	// }

	points := [][]float64{
		[]float64{1.2},
		[]float64{0.21},
		[]float64{2.21},
		[]float64{-3.21},
		[]float64{-4.21},

		[]float64{50.20},
		[]float64{63.21},
		[]float64{78.13},
	}

	fourier.GMM(points, 2)

	// pf, _ := os.Create("/tmp/points.csv")
	// cf, _ := os.Create("/tmp/centroids.csv")
	// defer pf.Close()
	// defer cf.Close()

	// fmt.Fprintf(pf, "X, Y, Z\n")
	// fmt.Fprintf(cf, "X, Y, Z\n")

	// for i, ms := range ms {
	// 	fmt.Fprintf(pf, "%f, %f, %d\n", points[i][0], points[i][1], ms.GetCluster())
	// }

	// for i, cc := range c {
	// 	fmt.Fprintf(cf, "%f, %f, %d\n", cc[0], cc[1], i)
	// }
}
