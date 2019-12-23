package main

import (
	"fmt"
	"gonum.org/v1/gonum/blas/blas32"
	//"gonum.org/v1/netlib/blas/netlib"
//	"gonum.org/v1/gonum/floats"
	"math"
	"math/rand"
	"sync"
	"time"
)

const vecLen = 256
const tasks = 1e6
const runs = 1e3
const threads = 8
const threshold = 0.4
const workFactor = tasks / threads
const totalDotProducts = runs * tasks

func toHuman(val float32, base10 bool) string {
	sufs := [5]string{"", "K", "M", "G", "T"}
	i := val
	idx := 0
	var base float32 = 1024
	if base10 {
		base = 1000
	}

	for ; i > base; i /= base {
		idx++
	}
	return fmt.Sprintf("%.2f%s", i, sufs[idx])
}

func initDB(container []blas32.Vector) {
	for i := 0; i < tasks; i++ {
		container[i] = getVector()
	}
}

func linearSearch(vect *blas32.Vector, db *[tasks]blas32.Vector, results *[tasks]float32, from int, to int, group *sync.WaitGroup) {
	for i := from; i < to; i++ {
		results[i] = blas32.Dot(*vect, db[i])
	}
	group.Done()
}

func getVector() (ret blas32.Vector) {
	var arr [vecLen]float32
	for i := 0; i < vecLen; i++ {
		arr[i] = rand.Float32()
	}
	vect := blas32.Vector{Inc: 1, Data: arr[:], N: vecLen}
	norm := blas32.Nrm2(vect)
	for i := 0; i < vecLen; i++ {
	    vect.Data[i] /= norm
	}
	return vect
}

func passesThreshold(score float32) bool {
	return score > threshold
}

func doSearch(candidate *blas32.Vector, container *[tasks]blas32.Vector, results *[tasks]float32) {
	var group sync.WaitGroup
	for wid := 0; wid < threads; wid++ {

		from := wid * workFactor
		to := from + workFactor

		group.Add(1)
		go linearSearch(candidate, container, results, from, to, &group)
	}

	group.Wait()
}

func filterResults(results [tasks]float32) (matches []int) {
//	matches, err := floats.Find(matches, passesThreshold, results[:], -1)

//	if err != nil {
//		panic(err)
//	}

	return matches
}

func main() {

	//blas32.Use(netlib.Implementation{})

	fmt.Println("Initializing vectors DB (this might take a while) for", toHuman(tasks, true), "vectors")
	fmt.Println("Going to use", toHuman(4*256*tasks, false)+"B", "Memory for DB")
	fmt.Println("Going to use", toHuman(4*tasks, false)+"B", "Memory for matches array")
	fmt.Println("Total RAM not counting program and GC overheads:", toHuman(4*(256+1)*tasks, false)+"B")

	if math.Mod(tasks, threads) != 0 {
		panic("tasks % threads should be 0")
	}

	start := time.Now()
	var container [tasks]blas32.Vector
	initDB(container[:])
	end := time.Now()
	took := end.Sub(start)
	fmt.Println("Initialized DB (took", took, "total and", took/tasks, "per vector)")
	fmt.Println()
	fmt.Printf("Running %s runs of %s tasks (total %s dot products) using %d threads\n",
		toHuman(runs, true), toHuman(tasks, true), toHuman(totalDotProducts, true), threads)

	var results [tasks]float32
	candidate := getVector()

	start = time.Now()
	var searchTook time.Duration
	var filteringTook time.Duration
	var totalMatches int32

	for i := 0; i < runs; i++ {
		if i%10 == 0 {
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
		totalMatches += int32(len(filterResults(results)))
		filteringTook += time.Now().Sub(temp)
	}
	end = time.Now()
	totalTook := end.Sub(start)

	dotProductsPerSecond := float32(totalDotProducts/time.Duration.Seconds(searchTook))
	flops := dotProductsPerSecond * vecLen
	fmt.Println()
	fmt.Println()
	fmt.Println("\rdone calculating", toHuman(totalDotProducts, true), "dot products in:", searchTook)
	fmt.Println()
	fmt.Println("Took ~", searchTook/totalDotProducts, "per dot product")
	fmt.Println("Calculated approx.", toHuman(dotProductsPerSecond, true)+"DP/s")
	fmt.Println("Calculation ran with approx. "+toHuman(flops, true)+"FLOPS")
	fmt.Println("Calculation accessed RAM with approx. "+toHuman(flops*4, true)+"B/s")
	fmt.Println("First 4 results from last run:", results[:4])
	fmt.Println()
	fmt.Println("Found a total of", toHuman(float32(totalMatches), true), "threshold-passing results")
	fmt.Println("Found on average", toHuman(float32(totalMatches/runs), true), "threshold-passing results per run")
	fmt.Println("Filteration took ~", filteringTook, "total")
	fmt.Println("Took ~", filteringTook/runs, "per run")
	fmt.Println()
	fmt.Println("Total time elapsed:", totalTook)
	fmt.Println("Total time per vector in DB:", totalTook/totalDotProducts)
}
