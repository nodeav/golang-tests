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
const VectSizeOnDisk = 256 * 8 + 4 + 4

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

func deserializeVector(buf []byte) (ret blas64.Vector) {
	if len(buf) < VectSizeOnDisk {
		panic("deserializeVector requires a specific size of buffer")
	}
	//fmt.Println("got buf", buf[:4], "with len", len(buf))
	size := binary.LittleEndian.Uint32(buf[0:4])
	//fmt.Println("using size", size)
	ret.Data = make([]float64, size)
	ret.Inc = int(binary.LittleEndian.Uint32(buf[4:8]))
	idx := 0
	for i := 8; i < VectSizeOnDisk; i+=8 {
		from := i
		to := i+8
		bits := binary.LittleEndian.Uint64(buf[from:to])
		float := math.Float64frombits(bits)
		ret.Data[idx] = float
		idx++
	}
	return
}

func loadFile(path string) (res []blas64.Vector) {
	f, err := os.Open(path)
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		panic(err)
	}
	numVecs := info.Size() / VectSizeOnDisk
	fmt.Println("Going to load", numVecs, "vectors!")
	res = make([]blas64.Vector, numVecs)

	buf := make([]byte, 64*VectSizeOnDisk)
	buf2 := make([]byte, 0, 2*64*VectSizeOnDisk)

	reader := bufio.NewReader(f)
	for true {
		n, err := reader.Read(buf)
		buf2 = append(buf2, buf...)
		for i := 0 ;; i++ {
			from := i * VectSizeOnDisk
			to := from + VectSizeOnDisk
			maybeVect := buf2[from:to]
			if n >= VectSizeOnDisk {
				//fmt.Println("Sending buf", maybeVect[0:4], "with len", len(maybeVect))
				res = append(res, deserializeVector(maybeVect))
				n -= VectSizeOnDisk
			} else {
				rest := buf[from:from+n]
				//fmt.Println("Copying", rest, "with len", len(rest))
				copy(buf2, rest)
				break
			}
		}
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}

	return
}

func (b Basic) Load(basePath string, vecs []blas64.Vector) {
	for _, path := range Readdir(basePath) {
		vecs = append(vecs, loadFile(path)...)
	}
}

func (b Basic) LoadChan(basePath string, vecs <-chan[]blas64.Vector) {

}
