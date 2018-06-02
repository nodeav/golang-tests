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
const tasks = 1e5
const runs = 1e3
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

func filterResults(results [tasks]float64) (matches []int) {
	matches, err := floats.Find(matches, passesThreshold, results[:], -1)

	if (err != nil) {
		panic(err)
	}

	return matches
}

func main() {

	fmt.Println("Initializing vectors DB (this might take a while) for", toHuman(tasks, true), "vectors")
	fmt.Println("Going to use", toHuman(8 * 256 * tasks, false)+"B", "Memory for DB")
	fmt.Println("Going to use", toHuman(8 * tasks, false)+"B", "Memory for matches array")
	fmt.Println("Total RAM not counting program and GC overheads:", toHuman(8 * (256 + 1) * tasks, false)+"B")

	if math.Mod(tasks, threads) != 0 {
		panic("tasks % threads should be 0")
	}

	start := time.Now()
	var container [tasks]blas64.Vector
	initDB(container[:])
	end := time.Now()
	took := end.Sub(start)
	fmt.Println("Initialized DB (took", took, "total and", took / tasks, "per vector)")
	fmt.Println()
	fmt.Printf("Running %s runs of %s tasks (total %s dot products) using %d threads\n",
		toHuman(runs, true), toHuman(tasks, true), toHuman(totalDotProducts, true), threads)

	var results [tasks]float64
	candidate := getVector()

	start = time.Now()
	var searchTook time.Duration
	var filteringTook time.Duration
	var totalMatches int64

	for i := 0; i < runs; i++ {
		if i % 10 == 0 {
			fmt.Printf("Running iter #%d\r", i)
		}
		if i+1 == runs {
			fmt.Print("<----- Done ----->")
		}
		// Measure searching
		temp := time.Now()
		doSearch(&candidate, &container, &results)
		searchTook += time.Now().Sub(temp)

		// Measure filtering
		temp = time.Now()
		totalMatches += int64(len(filterResults(results)))
		filteringTook += time.Now().Sub(temp)
	}
	end = time.Now()
	totalTook := end.Sub(start)

	fmt.Println()
	fmt.Println()
	fmt.Println("\rdone calculating", toHuman(totalDotProducts, true), "dot products in:", searchTook)
	fmt.Println()
	fmt.Println("Took ~", searchTook/totalDotProducts, "per dot product")
	fmt.Println("First 4 results from last run:", results[:4])
	fmt.Println()
	fmt.Println("Found a total of", toHuman(float64(totalMatches), true), "threshold-passing results")
	fmt.Println("Found on average", toHuman(float64(totalMatches/runs), true), "threshold-passing results per run")
	fmt.Println("Filteration took ~", totalTook, "total")
	fmt.Println("Took ~", totalTook/runs, "per run")
	fmt.Println()
	fmt.Println("Total time elapsed:", totalTook)
	fmt.Println("Total time per vector in DB:", totalTook / totalDotProducts)
}
