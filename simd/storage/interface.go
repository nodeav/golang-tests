package storage

import (
	"gonum.org/v1/gonum/blas/blas64"
)

type Storage interface {
	Store(basePath string, vecs []blas64.Vector) error
	Load(basePath string, vecs []blas64.Vector) error
	LoadChan(basePath string, vecs <-chan []blas64.Vector) error
}
