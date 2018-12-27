package fourier

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// MfccClusterisable - because every group of mfccs coefficient belongs to exacly one cluster
// we store the cluster Id along with the coefficients
type MfccClusterisable struct {
	coefficients []float64
	clusterID    int32
}

// Kmeans takes all the mfccs for several files
// which are of size nx39 (where n is file_len / 10ms)
// and returns k arrays of indices, where in array i we have the indices in the i-th cluster
func Kmeans(mfccsFloats [][]float64, k int) [][]int {
	variances := getCovarianceMatrixDiagonal(mfccsFloats)

	mfccs := make([]MfccClusterisable, len(mfccsFloats), len(mfccsFloats))
	for i, mfcc := range mfccsFloats {
		mfccs[i] = MfccClusterisable{
			coefficients: mfcc,
			clusterID:    -1,
		}
	}

	centroidIndices := make(map[int32]struct{})

	//First choose randomly the firts centroids
	// Keep generating a new random number until there are k keys in centroidIndices
	rand.Seed(time.Now().UTC().UnixNano())
	for len(centroidIndices) < k {
		ind := rand.Int31n(int32(len(mfccs)))
		if _, ok := centroidIndices[ind]; !ok {
			centroidIndices[ind] = struct{}{}
		}
	}

	centroids := make([][]float64, 0, k)
	for i := range centroidIndices {
		centroids = append(centroids, mfccs[i].coefficients)
	}

	iterations := 5
	// Group the documents in clusters and recalculate the new centroid of the cluster
	for times := 0; times < iterations; times++ {
		fmt.Printf("Iteration: %d\n", times)
		for i := range mfccs {
			mfccs[i].clusterID = findClosestCentroid(centroids, mfccs[i].coefficients, variances)
			fmt.Printf("%d %v\n", mfccs[i].clusterID, mfccs[i].coefficients)
		}

		centroids = findNewCentroids(mfccs, k)
	}

	fmt.Printf("Centroids:\n")
	for i, c := range centroids {
		fmt.Printf("%d %v\n", i, c)
	}

	return nil
}

func getCovarianceMatrixDiagonal(mfccs [][]float64) []float64 {
	variances := make([]float64, len(mfccs[0]), len(mfccs[0]))

	expectation := make([]float64, len(mfccs[0]), len(mfccs[0]))
	expectationSquared := make([]float64, len(mfccs[0]), len(mfccs[0]))
	for j := 0; j < len(mfccs[0]); j++ {
		for i := 0; i < len(mfccs); i++ {
			expectation[j] += mfccs[i][j]
			expectationSquared[j] += expectation[j] * expectation[j]
		}
	}

	for j := 0; j < len(mfccs[0]); j++ {
		variances[j] = expectationSquared[j]/float64(len(mfccs)) - expectation[j]*expectation[j]/float64(len(mfccs))
		if variances[j] < EPS {
			panic(fmt.Sprintf("%f", variances[j]))
		}
	}

	return variances
}

func findNewCentroids(mfccs []MfccClusterisable, k int) [][]float64 {
	centroids := make([][]float64, k, k)
	for i := range centroids {
		centroids[i] = make([]float64, len(mfccs[0].coefficients), len(mfccs[0].coefficients))
	}
	mfccsInCluster := make([]int, k, k)

	for _, mfcc := range mfccs {
		mfccsInCluster[mfcc.clusterID]++
		add(&centroids[mfcc.clusterID], mfcc.coefficients)
	}

	for i := range centroids {
		divide(&centroids[i], mfccsInCluster[i])
	}

	return centroids
}

func divide(x *[]float64, n int) {
	for i := 0; i < len(*x); i++ {
		(*x)[i] = (*x)[i] / float64(n)
	}
}

func add(x *[]float64, y []float64) {
	for i := 0; i < len(y); i++ {
		(*x)[i] += y[i]
	}
}

func findClosestCentroid(centroids [][]float64, mfcc []float64, variances []float64) int32 {
	// Returns positive infty if argument is >=0
	min := math.Inf(42)
	argmin := int32(-1)

	for i, centroid := range centroids {
		currentDistance := mahalanobisDistance(centroid, mfcc, variances)
		if currentDistance < min {
			min = currentDistance
			argmin = int32(i)
		}
	}
	return argmin
}

func mahalanobisDistance(x []float64, y []float64, variances []float64) float64 {
	sum := 0.0
	for i := 0; i < len(x); i++ {
		sum += (x[i] - y[i]) * (x[i] - y[i]) / variances[i]
	}

	return sum
}
