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
const tasks = 1e5
const threads = 8
const workFactor = tasks / threads

func initDB(container *[tasks][]float64) {
	for i := 0; i < tasks; i++ {
		container[i] = getVector()
	}
}

func linearSearch(vect *[]float64, db *[tasks][]float64, results *[tasks]float64, from int, to int, group *sync.WaitGroup) {
	for i := range db[from:to] {
		results[i] = floats.Dot(db[i], *vect)
	}
	group.Done()
}

func getVector() (ret []float64) {
	for i := 0; i < vecLen; i++ {
		ret = append(ret, rand.Float64())
	}
	return ret
}

func main() {

	if math.Mod(tasks, threads) != 0 {
		panic("times % maxWid should be 0")
	}

	start := time.Now()
	var container [tasks][]float64
	initDB(&container)
	end := time.Now()

	fmt.Println("initDB with", len(container), "elements took", end.Sub(start))
	fmt.Printf("\nRunning %.2e tasks using %d threads\n", tasks, threads)

	var results [tasks]float64
	candidate := getVector()

	start = time.Now()

	fmt.Println(start)

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

	fmt.Println(end)
	fmt.Println("done in:", took)
	fmt.Println("took ~", took/tasks, "per iteration")
	fmt.Println("Some results:", results[0:8])

}
