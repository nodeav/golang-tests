package storage
import (
	"os"
	"path/filepath"
	"gonum.org/v1/gonum/blas/blas64"
	"sync"
	"time"
	"fmt"
	"encoding/binary"
	"io"
	"math"
	"bufio"
)

type Basic struct {
	currfile 		*os.File
	currfilename	string
	vecsInFile 		int
	lock 			sync.Mutex
	Storage
}

const dbExt = ".vecs"

func WriteVect(w io.Writer, v *blas64.Vector) {
	length := len(v.Data)
	size := make([]byte, 4)
	inc := make([]byte, 4)
	data := make([]byte, 8 * length)


	binary.LittleEndian.PutUint32(size, uint32(length))
	binary.LittleEndian.PutUint32(inc, uint32(v.Inc))

	for i := 0; i < length; i++ {
		bits := math.Float64bits(v.Data[i])
		binary.LittleEndian.PutUint64(data[i*8:(i+1)*8], bits)
	}

	w.Write(size)
	w.Write(inc)
	_, err := w.Write(data)
	if err != nil {
		panic(err)
	}
}

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
			ret = append(ret, filepath.Join(basePath, f.Name()))
		}
	}
	return
}

func (b *Basic) getNewFile(filename string, basePath string) (error) {
	newFilename := filepath.Join(basePath, filename)
	flags := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	file, err := os.OpenFile(newFilename, flags, 0755)

	b.currfile = file
	b.currfilename = filename
	b.vecsInFile = 0
	return err
}

func (b *Basic) Store(basePath string, vecs []blas64.Vector) {

	filename := time.Now().Format("01-01-2018") + dbExt
	if b.currfilename == "" || b.currfilename != filename {
		b.lock.Lock()
		err := b.getNewFile(filename, basePath)
		b.lock.Unlock()
		if err != nil {
			panic(err)
		}
	}
	writer := bufio.NewWriter(b.currfile)
	for i := range vecs {
		WriteVect(writer, &(vecs[i]))
	}
	b.vecsInFile += len(vecs)
	err := writer.Flush()
	if err != nil {
		panic(err)
	}
	b.currfile.Sync()
	fmt.Println("vecsInFile", b.vecsInFile)
}

func (b Basic) Load(basePath string, vecs []blas64.Vector) {
	Readdir(basePath)
}

func (b Basic) LoadChan(basePath string, vecs <-chan[]blas64.Vector) {

}
