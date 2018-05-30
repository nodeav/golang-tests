package main

import (
	"fmt"
	"gonum.org/v1/gonum/floats"
	"math"
	"math/rand"
	"sync"
	"time"
)

const vecLen = 256
const tasks = 1e6
const threshold = 0.35
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

func initDB(container *[tasks][]float64) {
	for i := 0; i < tasks; i++ {
		container[i] = getVector()
	}
}

func linearSearch(vect *[]float64, db *[tasks][]float64, results *[tasks]float64, from int, to int, group *sync.WaitGroup) {
	for i := from; i < to; i++ {
		results[i] = floats.Dot(db[i], *vect)
	}
	group.Done()
}

func getVector() (ret []float64) {
	var arr [vecLen]float64
	for i := 0; i < vecLen; i++ {
		arr[i] = rand.Float64()
	}
	ret = arr[:]
	floats.Mul(ret, ret)
	scaleFactor := 1 / math.Sqrt(floats.Sum(ret))
	floats.Scale(scaleFactor, ret)
	return ret
}

func passesThreshold(score float64) bool {
	return score > threshold
}

func main() {

	if math.Mod(tasks, threads) != 0 {
		panic("times % maxWid should be 0")
	}

	start := time.Now()
	var container [tasks][]float64
	initDB(&container)
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

	fmt.Println("Found", len(matches), "threshold-passing results in", took)
	fmt.Println("Took ~", took/tasks, "per count query")
	fmt.Println("Some match indices:", matches[0:8], "\n")

}
