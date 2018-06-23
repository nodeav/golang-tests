package main

import (
	"sync"
	DB "simd/simd/search"
	"fmt"
	"time"
)

func main() {
	const amountCands = 256
	const dbLen = 1.5e6
	const workers = 8
	const workStep = dbLen/workers

	db := &DB.Linear{Size: dbLen}
	cands := &DB.Linear{Size: amountCands}

	start := time.Now()
	db.InitRandom()
	cands.InitRandom()
	end := time.Now()
	initTook := end.Sub(start)
	fmt.Println("init took", initTook)
	g := sync.WaitGroup{}
	g.Add(workers)
	results := make([][]int, amountCands)
	for i := range results {
		results[i] = make([]int, dbLen)
	}

	start = time.Now()
	for i := 0; i < workers; i++ {
		from := i * workStep
		to := from + workStep
		go db.SearchWaitGroup(cands.Vectors, 0.55, &g, from, to, results)
	}

	g.Wait()
	end = time.Now()
	took := end.Sub(start)
	tookPer := took / (dbLen*amountCands)
	start = time.Now()
	var nResults int
	for i := range results {
		for _, x := range results[i] {
			if x != 0 {
				nResults++
			}
		}
	}
	end = time.Now()
	ftook := end.Sub(start)
	numSearches := amountCands*dbLen
	fmt.Println("got numResults:", nResults, "out of", numSearches,"\ntook:", took, "\ntook per vec", tookPer, "\nfilterting took", ftook)
}
