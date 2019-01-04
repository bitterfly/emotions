package fourier

import (
	"fmt"
	"math"
)

func GMM(mfccsFloats [][]float64, k int) {
	X, expectations, variances, numInCluster := KMeans(mfccsFloats, k)
	em(expectations, variances, numInCluster, X, k)
}

func em(expectations [][]float64, variances [][]float64, numInCluster []int, X []MfccClusterisable, k int) {
	// start with k-mean clusters expectations, variances and we use the number of points in a cluster for phi
	// and choose

	phi := make([]float64, k, k)
	for i := 0; i < k; i++ {
		phi[i] = float64(numInCluster[i]) / float64(len(X))
	}

	// weights

	fmt.Printf("X: %v\n", X)

	//compute the weights
	for step := 0; step < 10; step++ {
		fmt.Printf("Step: %d, Log likelihood: %f\n", step, logLikelihood(X, phi, expectations, variances, k))
		fmt.Printf("phi: %v\n", phi)
		fmt.Printf("expectations: %v\n", expectations)
		fmt.Printf("variances: %v\n", variances)

		fmt.Printf("len(X): %d\n", len(X))
		w := make([][]float64, len(X), len(X))
		for i := 0; i < len(X); i++ {
			sum := 0.0

			for j := 0; j < k; j++ {
				w[i] = make([]float64, k, k)
				w[i][j] = phi[j] * N(X[i].coefficients, expectations[j], variances[j])
				fmt.Printf("w[%d][%d]= %f,phi: %f\n", i, j, w[i][j], phi[j])
				fmt.Printf("N: %f,c: %v, exp: %v, var: %v\n", N(X[i].coefficients, expectations[j], variances[j]), X[i].coefficients, expectations[j], variances[j])
				sum += w[i][j]
			}
			fmt.Printf("sum: %f\n", sum)
			divide(&w[i], sum)
			fmt.Printf("w[%d] = %v\n", i, w[i])
		}

		N := make([]float64, k, k)
		for i := 0; i < len(X); i++ {
			for j := 0; j < k; j++ {
				N[j] += w[i][j]
			}
		}

		fmt.Printf("weights: %v\n", w)
		fmt.Printf("Ns: %v\n", N)
		// compute the μ Σ ɸ
		for j := 0; j < k; j++ {
			zero(&(expectations[j]))
			zero(&(variances[j]))
		}
		zero(&phi)

		for i := 0; i < len(X); i++ {
			for j := 0; j < k; j++ {
				add(&expectations[j], multiplied(X[i].coefficients, w[i][j]))

				fmt.Printf("Expectation[%d]: %f\n", j, expectations[j])
				diagonal := minused(X[i].coefficients, expectations[j])
				square(&diagonal)

				add(&variances[j], multiplied(diagonal, w[i][j]))
			}
		}

		for j := 0; j < k; j++ {
			divide(&(expectations[j]), N[j])
			divide(&(variances[j]), N[j])
			phi[j] = N[j] / float64(len(X))
		}
	}

}

func getDeterminant(variance []float64) float64 {
	det := 1.0
	for i := 0; i < len(variance); i++ {
		det *= variance[i]
	}
	return det
}

func N(xi []float64, expectation []float64, variance []float64) float64 {
	var exp float64
	for i := 0; i < len(xi); i++ {
		exp += (xi[i] - expectation[i]) * (xi[i] - expectation[i]) / variance[i]
	}

	fmt.Printf("exp: %f\n", exp)

	//N: 0.000000,c: [6.21], exp: [0.705], var: [0.24502500000000005]

	return math.Exp(-0.5*exp) / math.Sqrt(math.Pow(2.0*math.Pi, float64(len(xi)))*getDeterminant(variance))
}

func logLikelihood(X []MfccClusterisable, phi []float64, expectations [][]float64, variances [][]float64, k int) float64 {
	sum := 0.0
	for i := 0; i < len(X); i++ {
		for j := 0; j < k; j++ {
			fmt.Printf("%f N(%v, %v, %v) - %.8f\n", phi[j], X[i].coefficients, expectations[j], variances[j], phi[j]*N(X[i].coefficients, expectations[j], variances[j]))
			sum += math.Log(phi[j] * N(X[i].coefficients, expectations[j], variances[j]))
		}
	}
	return sum
}
