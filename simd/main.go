package main

import (
	"fmt"
	"gonum.org/v1/gonum/blas/blas64"
	DB "simd/simd/search"
	"simd/simd/storage"
	"time"
	"flag"
	"os/exec"
)

func purge() {
	purgeStart := time.Now()
	cmd := exec.Command("purge")
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println("purge took", time.Now().Sub(purgeStart))
	time.Sleep(1 * time.Second)
}

func saveDB(dbPath string, dbSize int) {
	var db []blas64.Vector

	initRandomStart := time.Now()
	DB.InitRandom(&db, dbSize)
	initRandomTook := time.Now().Sub(initRandomStart)
	fmt.Println("Initializing a random DB took", initRandomTook)
	fmt.Printf("< save > Some of db: %+v\n", db[0].Data[:4])

	store := &storage.Basic{}
	start := time.Now()
	store.Store(dbPath, db)
	end := time.Now()
	fsave := end.Sub(start)
	fmt.Println("Save took:", fsave)
}

func loadDB(dbPath string) {
	fmt.Println("DB directory contains these files:", storage.Readdir(dbPath))

	var db []blas64.Vector
	s := storage.Basic{}
	loadStart := time.Now()
	s.Load(dbPath, &db)
	loadEnd := time.Now()
	fmt.Printf("< load > Some of db: %+v\n", db[0].Data[:4])
	loadTook := loadEnd.Sub(loadStart)
	fmt.Println("Load took:", loadTook)

	numNoInc := 0
	for i := 0; i < len(db); i++ {
		if db[i].Inc != 1 {
			numNoInc++
		}
	}
	fmt.Printf("Found [ %d ] entries of which [ %d ] were invalid\n", len(db)-numNoInc, numNoInc)
}

func main() {
	var load, save bool
	var dbPath string
	var dbSize int

	flag.BoolVar(&load, "load", true, "Benchmark loading the db")
	flag.BoolVar(&save, "save", true, "Benchmark saving the db")
	flag.StringVar(&dbPath,"db", "./db", "Path to load/save the db. (Should exist)")
	flag.IntVar(&dbSize, "size", 5e5, "Amount of db entries")
	flag.Parse()

	flag.VisitAll(func (f *flag.Flag) {
		fmt.Println("Key", f.Name, "Value", f.Value)
	})

	if save {
		fmt.Println("save is", save)
		purge()
		saveDB(dbPath, dbSize)
	}

	if load {
		fmt.Println("load is", load)
		purge()
		loadDB(dbPath)
	}
}
