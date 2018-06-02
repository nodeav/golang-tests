package main

import (
	"fmt"
	//"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/blas/blas64"
	//"gonum.org/v1/netlib/blas/netlib"
	"math"
	"math/rand"
	"sync"
	"time"
)

const vecLen = 256
const tasks = 5e5
const runs = 1e4
const threads = 8
const threshold = 0.4
const workFactor = tasks / threads
const totalDotProducts = runs * tasks

func toHuman(val float64, base10 bool) string {
	sufs := [5]string {"", "K", "M", "G", "T"}
	i := val
	idx := 0
	var base float64 = 1024
	if base10 {base = 1000}

	for ; i > base; i /= base {
		idx++
	}
	return fmt.Sprintf("%.2f%s", i, sufs[idx])
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

func doSearch(candidate *blas64.Vector, container *[tasks]blas64.Vector, results *[tasks]float64) {
	var group sync.WaitGroup
	for wid := 0; wid < threads; wid++ {

		from := wid * workFactor
		to := from + workFactor

		group.Add(1)
		go linearSearch(candidate, container, results, from, to, &group)
	}

	group.Wait()
}

func main() {

	//blas64.Use(netlib.Implementation{})

	if math.Mod(tasks, threads) != 0 {
		panic("tasks % threads should be 0")
	}

	start := time.Now()
	var container [tasks]blas64.Vector
	initDB(container[:])
	end := time.Now()

	fmt.Println("\ninitDB with", len(container), "elements took", end.Sub(start), "and",
		toHuman(8 * 256 * tasks, false), "Memory")

	fmt.Printf("Running %s runs of %s tasks (total %s dot products) using %d threads\n",
		toHuman(runs, true), toHuman(tasks, true), toHuman(totalDotProducts, true), threads)

	var results [tasks]float64
	candidate := getVector()

	start = time.Now()
	for i := 0; i < runs; i++ {
		if i % 10 == 0 || i+1 == runs {
			fmt.Printf("\rRunning iter #%d", i)
		}
		doSearch(&candidate, &container, &results)
	}
	end = time.Now()
	took := end.Sub(start)

	fmt.Println("\rdone calculating", toHuman(totalDotProducts, true), "dot products in:", took)
	fmt.Println("took ~", took/(totalDotProducts), "per multiplication")
	fmt.Println("Some results from last run:", results[:8])

	//start = time.Now()
	//var matches []int
	//matches, err := floats.Find(matches, passesThreshold, results[:], -1)
	//
	//if (err != nil) {
	//	panic(err)
	//}
	//
	//end = time.Now()
	//took = end.Sub(start)
	//
	//someMatches := matches[0:8]
	//fmt.Println("Found", len(matches), "threshold-passing results in", took)
	//fmt.Println("Took ~", took/tasks, "per count query")
	//fmt.Println("Some match indices:", someMatches, "\n")
	//fmt.Print("Some match scores: ")
	//for _, k := range someMatches {
	//	fmt.Print(k, ": ", results[k], ", ")
	//}
	//fmt.Println()

}
