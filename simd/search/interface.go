package search

import 	(
	"gonum.org/v1/gonum/blas/blas64"
)


type Searcher interface {
	Float64(db, candidates []blas64.Vector, threshold float64, from, to int, results [][]int)
	Chan(db, candidates []blas64.Vector, threshold float64, from, to int, results chan[]float64)
	WaitGroupF64(db, candidates []blas64.Vector, threshold float64, from, to int, results [][]int)
}