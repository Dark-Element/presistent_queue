package services

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
)

func NewFileQueue(prefix string, maxFileSize int64) *FileQueue {
	f := createFile(prefix, 0)
	r := make(chan *os.File)
	cr := createReader(f)

	fq := FileQueue{w: f, r: r, prefix: prefix, maxFileSize: maxFileSize, currentR: bufio.NewReader(cr), currentRW: cr, currentFilename: cr.Name()}
	return &fq
}

type iFile interface {
	Push(data bytes.Buffer) bool
	Pop(c int64) bytes.Buffer
	Close()
}

type FileQueue struct {
	wMutex sync.Mutex
	rMutex sync.Mutex

	w               *os.File
	r               chan *os.File
	currentR        *bufio.Reader
	currentRW       *os.File
	currentFilename string

	prefix      string
	currentNum  int64
	maxFileSize int64
}

//push to file
//lock the request until flushed (ACID compliance)
func (f *FileQueue) Push(data bytes.Buffer) {
	f.wMutex.Lock()
	data.WriteString("\n")

	f.w.WriteString(data.String())

	fi, err := f.w.Stat()
	if err != nil {
		log.Panic("BBBBBBBBBBB")
	}
	if fi.Size() >= f.maxFileSize {
		f.currentNum++
		cw := f.w
		go func() {
			f.r <- createReader(cw)
			cw.Close()
		}()

		f.w = createFile(f.prefix, f.currentNum)
	}
	f.wMutex.Unlock()
}

//change the output to bytes.buffer channel in order to utilize less memory
func (f *FileQueue) Pop(c int) bytes.Buffer {
	f.rMutex.Lock()
	b := bytes.Buffer{}
	var line string
	var err error
	for i := 0; i < c; i++ {
		line, err = f.currentR.ReadString('\n')
		b.WriteString(line)
		if err != nil {
			fmt.Println(err)
			//delete file
			if f.w.Name() != f.currentRW.Name(){
				os.Remove(f.currentFilename)
			}
			//pull new one from channel
			select {
			case x, ok := <-f.r:
				if ok {
					fmt.Println("Pulled new file")
					f.currentRW = x
					f.currentR = bufio.NewReader(x)
					f.currentFilename = x.Name()
				}
			default:
				fmt.Println("No files to read from")
				break
			}
		}
	}
	f.rMutex.Unlock()
	return b
}

func (f *FileQueue) Close() {

}

func createFile(prefix string, lastNum int64) *os.File {
	b := bytes.Buffer{}
	cd, _ := os.Getwd()
	b.WriteString(cd)
	b.WriteString("/files/")
	b.WriteString(prefix)
	b.WriteString("-")
	b.WriteString(strconv.FormatInt(lastNum+1, 10))
	fmt.Println(b.String())
	f, err := os.OpenFile(b.String(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Panic(err)
	}
	return f
}

func createReader(f *os.File) *os.File{
	f, err := os.OpenFile(f.Name(), os.O_RDONLY, 0777)
	if err != nil {
		log.Panic(err)
	}
	return f
}
