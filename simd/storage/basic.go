package storage
import (
	"os"
	"path/filepath"
	"gonum.org/v1/gonum/blas/blas64"
)

type Basic struct {
	Storage
}

const dbExt = ".vecs"

func Readdir(basePath string) (ret []string) {
	f, err := os.Open(basePath)
	if err != nil {
		panic(err)
	}
	dir, err := f.Readdir(-1)
	if err != nil {
		panic(err)
	}
	for _, f := range dir {
		isMatch := filepath.Ext(f.Name()) == dbExt
		if err != nil {
			panic(err)
		}
		if isMatch && !f.IsDir() {
			ret = append(ret, f.Name())
		}
	}
	return
}

func (b Basic) Store(basePath string, vecs []blas64.Vector) {

}

func (b Basic) Load(basePath string, vecs []blas64.Vector) {
	Readdir(basePath)
}

func (b Basic) LoadChan(basePath string, vecs <-chan[]blas64.Vector) {

}
