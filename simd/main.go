package main

import (
	"fmt"
	"gonum.org/v1/gonum/blas/blas64"
	DB "simd/simd/search"
	"simd/simd/storage"
	"sync"
	"time"
)

func main() {
	const amountCands = 128
	const dbLen = 1e3
	const workers = 8
	const workStep = dbLen / workers

	fmt.Println("loaded files:", storage.Readdir("./db"))
	s := storage.Basic{}
	searcher := &DB.Linear{}

	var db, cands []blas64.Vector

	start := time.Now()
	DB.InitRandom(&cands, amountCands)
	//DB.InitRandom(&db, dbLen)
	s.Load("./db", &db)
	end := time.Now()
	initTook := end.Sub(start)
	fmt.Println("init took", initTook)

	//store := &storage.Basic{}
	//start = time.Now()
	//store.Store("./db", db)
	//end = time.Now()
	//fsave := end.Sub(start)
	//fmt.Println("File saving took:", fsave)

	g := sync.WaitGroup{}
	g.Add(workers)
	results := make([][]int, amountCands)
	for i := range results {
		results[i] = make([]int, dbLen)
	}

	numNoInc := 0
	for i := 0; i < len(db); i++ {
		if db[i].Inc != 1 {
			numNoInc++
		}
	}
	fmt.Printf("There are %d with inc==1 and %d without\n", len(db)-numNoInc, numNoInc)

	start = time.Now()
	for i := 0; i < workers; i++ {
		from := i * workStep
		to := from + workStep
		go searcher.WaitGroupF64(db, cands, 0.45, &g, from, to, results)
	}

	g.Wait()
	end = time.Now()
	took := end.Sub(start)
	tookPer := took / (dbLen * amountCands)
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
	numSearches := amountCands * dbLen
	fmt.Println("got numResults:", nResults, "out of", numSearches, "\ntook:", took, "\ntook per vec", tookPer, "\nfilterting took", ftook)
}
