package search

import (
	"gonum.org/v1/gonum/blas/blas64"
	"sync"
	"errors"
	"math/rand"
	"fmt"
)
type Linear struct  {
	Vectors []blas64.Vector
	Size    int
	DB
}

func (db Linear) SearchWaitGroup(
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
			result := blas64.Dot(1, candidates[j], db.Vectors[i])
			if result > threshold {
				results[j][i] = i
			}
		}
	}

	return
}

func shuf (slice []float64) (res []float64) {
	res = make([]float64, len(slice))
	copy(res, slice)
	rand.Shuffle(len(res), func(i, j int) {
		res[i], res[j] = res[j], res[i]
	})
	return
}

func (db *Linear) InitRandom() {
	const vecLen = 256
	vecData := make([]float64, 256)
	for j := 0; j < vecLen; j++ {
		vecData[j] = rand.Float64()
	}
	normVec := blas64.Vector{Inc: 1, Data: vecData}
	blas64.Nrm2(1, normVec)
	vecData = normVec.Data
	for i := 0; i < db.Size; i++ {
		shuffed := shuf(vecData)
		db.Vectors = append(db.Vectors, blas64.Vector{Inc: 1, Data: shuffed})
	}
	return
}
