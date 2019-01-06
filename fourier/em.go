package fourier

import (
	"fmt"
	"math"
	"os"
)

// GMM returns the k gaussian mixures for the given data
func GMM(mfccsFloats [][]float64, k int) ([]float64, [][]float64, [][]float64) {
	fmt.Printf("kMeans\n")
	X, expectations, variances, numInCluster := KMeans(mfccsFloats, k)
	fmt.Printf("NumInClusters: %v\n", numInCluster)
	fmt.Printf("EM\n")
	return em(expectations, variances, numInCluster, X, k)
}

func em(expectations [][]float64, variances [][]float64, numInCluster []int, X []MfccClusterisable, k int) ([]float64, [][]float64, [][]float64) {
	// start with k-mean clusters expectations, variances and we use the number of points in a cluster for phi
	// and choose

	phi := make([]float64, k, k)
	for i := 0; i < k; i++ {
		phi[i] = float64(numInCluster[i]) / float64(len(X))
	}

	// weights
	epsilon := 0.000001
	f, _ := os.Create("/tmp/ws")

	prevLikelihood := 0.0
	likelihood := 0.0
	for step := 0; step < 10; step++ {
		w := make([][]float64, len(X), len(X))
		maximums := make([]float64, len(X), len(X))
		for i := 0; i < len(maximums); i++ {
			maximums[i] = math.Inf(-1)
		}

		fmt.Fprintf(f, "Step %d\n", step)

		for i := 0; i < len(X); i++ {
			w[i] = make([]float64, k, k)

			for j := 0; j < k; j++ {
				w[i][j] = math.Log(phi[j]) + N(X[i].coefficients, expectations[j], variances[j])

				fmt.Fprintf(f, "w[%d][%d] = %f, phi[%d]: %f, N: %f\n", i, j, w[i][j], j, phi[j], N(X[i].coefficients, expectations[j], variances[j]))
				fmt.Fprintf(f, "X: %v\nexp: %v\nvar: %v\n\n", X[i].coefficients, expectations[j], variances[j])

				if maximums[i] < w[i][j] {
					maximums[i] = w[i][j]
				}
			}
			var sum float64
			for j := 0; j < k; j++ {
				if i == 10 {
				}
				if w[i][j] < maximums[i]-10 {
					w[i][j] = 0
				} else {
					w[i][j] = math.Exp(w[i][j] - maximums[i])
					sum += w[i][j]
				}
			}

			fmt.Fprintf(f, "sum[%d]: %f\n", i, sum)

			divide(&w[i], sum)
		}

		N := make([]float64, k, k)
		for i := 0; i < len(X); i++ {
			for j := 0; j < k; j++ {
				N[j] += w[i][j]
			}
		}

		for i := 0; i < k; i++ {
			if N[i] < epsilon {
				N[i] = epsilon
			}
		}

		fmt.Printf(fmt.Sprintf("Step[%d], w: %v\nNs: %v\nphi: %v\nexp: %v\nvar: %v\n\n", step, w, N, phi, expectations, variances))

		for j := 0; j < k; j++ {
			zero(&(expectations[j]))
			zero(&(variances[j]))
		}
		zero(&phi)

		for i := 0; i < len(X); i++ {
			for j := 0; j < k; j++ {
				add(&expectations[j], multiplied(X[i].coefficients, w[i][j]))

				// fmt.Printf("Expectation[%d]: %f\n", j, expectations[j])
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

		likelihood = logLikelihood(X, phi, expectations, variances, k)
		fmt.Printf("Log likelihood: %f\n", likelihood)

		if epsDistance(likelihood, prevLikelihood, 0.00001) {
			fmt.Printf("Break on step: %d\n", step)
			break
		}

		prevLikelihood = likelihood
	}

	return phi, expectations, variances
}

func epsDistance(a, b, e float64) bool {
	return (a-b < e && a-b > -e)
}

func getDeterminant(variance []float64) float64 {
	det := 1.0
	for i := 0; i < len(variance); i++ {
		det *= variance[i]
	}
	return det
}

// EvaluateVector returns the likelihood a given vector
func EvaluateVector(X []float64, phi []float64, expectations [][]float64, variances [][]float64, k int) float64 {
	return logLikelihoodFloat(X, phi, expectations, variances, k)
}

func N(xi []float64, expectation []float64, variance []float64) float64 {
	var exp float64
	for i := 0; i < len(xi); i++ {
		exp += (xi[i] - expectation[i]) * (xi[i] - expectation[i]) / variance[i]
	}

	// return math.Exp(-0.5*exp) / math.Sqrt(math.Pow(2.0*math.Pi, float64(len(xi)))*getDeterminant(variance))
	// return log of this
	return -0.5 * (exp - float64(len(xi))*math.Log(2*math.Pi) - math.Log(getDeterminant(variance)))
}

func logLikelihoodFloat(X []float64, phi []float64, expectations [][]float64, variances [][]float64, k int) float64 {
	sum := 0.0
	for j := 0; j < k; j++ {
		sum += math.Log(phi[j] * math.Exp(N(X, expectations[j], variances[j])))
	}
	return sum
}

func logLikelihood(X []MfccClusterisable, phi []float64, expectations [][]float64, variances [][]float64, k int) float64 {
	sum := 0.0
	for i := 0; i < len(X); i++ {
		for j := 0; j < k; j++ {
			sum += math.Log(phi[j] * math.Exp(N(X[i].coefficients, expectations[j], variances[j])))
		}
	}
	return sum
}
