package main

import (
	"flag"
	"fmt"
	"gonum.org/v1/gonum/blas/blas64"
	"log"
	"os"
	"os/exec"
	"runtime/pprof"
	DB "simd/simd/search"
	"simd/simd/storage"
	"time"
)

func toHuman(val float64, base10 bool) string {
	sufs := [5]string{"", "K", "M", "G", "T"}
	i := val
	idx := 0
	var base float64 = 1024
	if base10 {
		base = 1000
	}

	for ; i > base; i /= base {
		idx++
	}
	return fmt.Sprintf("%.2f%s", i, sufs[idx])
}

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
	vecsPerSecond := int64(len(db)) / int64(loadTook/time.Second)
	fmt.Printf("Vecs per second: %s/s\n", toHuman(float64(vecsPerSecond), false))
	fmt.Printf("Read speed: %sB/s\n", toHuman(float64(vecsPerSecond*storage.VectSizeOnDisk), false))
}

func main() {
	var load, save bool
	var dbPath, cpuprofilePath string
	var dbSize int

	flag.BoolVar(&load, "load", true, "Benchmark loading the db")
	flag.BoolVar(&save, "save", true, "Benchmark saving the db")
	flag.StringVar(&dbPath, "db", "./db", "Path to load/save the db. (Should exist)")
	flag.IntVar(&dbSize, "size", 5e5, "Amount of db entries")
	flag.StringVar(&cpuprofilePath, "cpuprofile", "./cpu-prof", "write cpu profile to `file`")

	flag.Parse()

	flag.VisitAll(func(f *flag.Flag) {
		fmt.Println("Key", f.Name, "Value", f.Value)
	})

	if cpuprofilePath != "" {
		f, err := os.Create(cpuprofilePath)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

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
