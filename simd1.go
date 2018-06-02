package main

import (
	"fmt"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/blas/blas64"
	"math"
	"math/rand"
	"sync"
	"time"
)

const vecLen = 256
const tasks = 2.5e5
const threshold = 0.43
const threads = 8
const workFactor = tasks / threads

func toHuman(val float64) string {
	sufs := [5]string {"", "K", "M", "G", "T"}
	i := val
	idx := 0
	for ; i > 1024; i /= 1024 {
		idx++
	}
	return fmt.Sprintf("%.2f%sB", i, sufs[idx])
}

func initDB(container []blas64.Vector) {
	for i := 0; i < tasks; i++ {
		container[i] = getVector()
	}
}

func linearSearch(vect *blas64.Vector, db *[tasks]blas64.Vector, results *[tasks]float64, from int, to int, group *sync.WaitGroup) {
	for i := from; i < to; i++ {
		results[i] = blas64.Dot(1, *vect, db[i])
	}
	group.Done()
}

func getVector() (ret blas64.Vector) {
	var arr [vecLen]float64
	for i := 0; i < vecLen; i++ {
		arr[i] = rand.Float64()
	}
	vect := blas64.Vector{Inc: 1, Data: arr[:]}
	blas64.Nrm2(1, vect)
	return vect
}

func passesThreshold(score float64) bool {
	return score > threshold
}

func main() {

	if math.Mod(tasks, threads) != 0 {
		panic("times % maxWid should be 0")
	}

	start := time.Now()
	var container [tasks]blas64.Vector
	initDB(container[:])
	end := time.Now()

	fmt.Println("\ninitDB with", len(container), "elements took", end.Sub(start), "and", toHuman(8 * 256 * tasks), "Memory")
	fmt.Printf("Running %.2e tasks using %d threads\n", tasks, threads)

	var results [tasks]float64
	candidate := getVector()

	start = time.Now()

	var group sync.WaitGroup
	for wid := 0; wid < threads; wid++ {

		from := wid * workFactor
		to := from + workFactor

		group.Add(1)
		go linearSearch(&candidate, &container, &results, from, to, &group)
	}

	group.Wait()
	end = time.Now()
	took := end.Sub(start)

	fmt.Println("done calculating dot products in:", took)
	fmt.Println("took ~", took/tasks, "per multiplication")
	fmt.Println("Some results:", results[:8])

	start = time.Now()
	var matches []int
	matches, err := floats.Find(matches, passesThreshold, results[:], -1)

	if (err != nil) {
		panic(err)
	}

	end = time.Now()
	took = end.Sub(start)

	someMatches := matches[0:8]
	fmt.Println("Found", len(matches), "threshold-passing results in", took)
	fmt.Println("Took ~", took/tasks, "per count query")
	fmt.Println("Some match indices:", someMatches, "\n")
	fmt.Print("Some match scores: ")
	for _, k := range someMatches {
		fmt.Print(k, ": ", results[k], ", ")
	}
	fmt.Println()

}
