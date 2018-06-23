package main

import (
	"sync"
	DB "simd/simd/search"
	"fmt"
	"time"
	"gonum.org/v1/gonum/blas/blas64"
)

func main() {
	const amountCands = 128
	const dbLen = 5e5
	const workers = 8
	const workStep = dbLen/workers

	searcher := &DB.Linear{}

	var db, cands []blas64.Vector

	start := time.Now()
	DB.InitRandom(&cands, amountCands)
	DB.InitRandom(&db, dbLen)
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
		go searcher.WaitGroupF64(db, cands, 0.45, &g, from, to, results)
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
