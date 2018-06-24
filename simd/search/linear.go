package search

import (
	"errors"
	"fmt"
	"gonum.org/v1/gonum/blas/blas64"
	"math/rand"
	"sync"
)

type Linear struct {
	Searcher
}

func (Linear) WaitGroupF64(
	vecs []blas64.Vector,
	candidates []blas64.Vector,
	threshold float64,
	group *sync.WaitGroup,
	from int,
	to int,
	results [][]int) (err error) {

	defer group.Done()
	if cap(results) < len(candidates) {
		str := fmt.Sprintf("results capacity is too low: %d < %d", cap(results), len(candidates))
		err = errors.New(str)
		return
	}

	for i := from; i < to; i++ {
		for j := 0; j < len(candidates); j++ {
			//fmt.Println("accessing vecs[i]:", len(vecs), i, "with from to", from, to)
			result := blas64.Dot(1, candidates[j], vecs[i])
			if result > threshold {
				results[j][i] = i
			}
		}
	}

	return
}

func shuf(slice []float64) (res []float64) {
	res = make([]float64, len(slice))
	copy(res, slice)
	shufInPlace(&res)
	return
}

func shufInPlace(slc *[]float64) {
	rand.Shuffle(len(*slc), func(i, j int) {
		(*slc)[i], (*slc)[j] = (*slc)[j], (*slc)[i]
	})
}

func InitRandom(container *[]blas64.Vector, sz int) {
	const vecLen = 256

	vecData := make([]float64, 256)
	*container = make([]blas64.Vector, sz)
	for j := 0; j < vecLen; j++ {
		vecData[j] = rand.Float64()
	}
	normVec := blas64.Vector{Inc: 1, Data: vecData}
	blas64.Nrm2(1, normVec)
	vecData = normVec.Data

	for i := range *container {
		shufInPlace(&vecData)
		(*container)[i] = blas64.Vector{Inc: 1, Data: vecData}
	}
	return
}
