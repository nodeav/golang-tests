package main

import (
	"fmt"
	"math"
	"gonum.org/v1/gonum/floats"
	"math/rand"
	"time"
	"sync"
)

const vecLen = 256
const tasks = 1e9
const threads = 8
const workFactor = tasks / threads

func main() {
	fmt.Printf("Running %.0f tasks using %d threads\n", tasks, threads)

	if math.Mod(tasks, threads) != 0 {
		panic("times % maxWid should be 0")
	}

	var vect1 	[]float64
	var vect2 	[]float64
	var results [tasks]float64

	for i := 0; i < vecLen; i++ {
		vect1 = append(vect1, rand.Float64())
		vect2 = append(vect2, rand.Float64())
	}

	start := time.Now()

	fmt.Println(start)

	var group sync.WaitGroup
	for wid := 0; wid < threads; wid++ {

		from :=  wid * workFactor
		to 	 := from + workFactor

		group.Add(1)
		go (func(from int, to int, group *sync.WaitGroup, results *[tasks]float64) {
			for i := from; i < to; i++ {
				results[i] = floats.Dot(vect1, vect2)
			}
			group.Done()
		})(from, to, &group, &results)
	}

	group.Wait()
	end := time.Now()
	took := end.Sub(start)

	fmt.Println(end)
	fmt.Println("done in:", took)
	fmt.Println("took ~", took / tasks, "per iteration")
	fmt.Println("resulting multiplication:", results[15000000])

}
