package search

import 	(
	"gonum.org/v1/gonum/blas/blas64"
)


type DB interface {
	SearchFloat64(candidates []blas64.Vector, threshold float64, from, to, results []float64)
	SearchChan(candidates []blas64.Vector, threshold float64, from, to, results chan[]float64)
}