package main

import (
	"fmt"
	"gonum.org/v1/gonum/floats"
	"github.com/slimsag/rand/simd"
)

func main() {
	t := [8]float64 { 0, 1, 2, 3, 4, 5, 6, 7 }
	a := simd.Vec64 { t[0], t[1], t[2], t[3] }
	b := simd.Vec64 { t[4], t[5], t[6], t[7] }
	c := simd.Vec64Mul(a, b)
	e := floats.Dot(t[0:4], t[4:8])
	fmt.Println("simd multiplication result: ", c)
	fmt.Println("gonum floats multiplication result", e)
}
