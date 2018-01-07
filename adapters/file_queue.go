package adapters

import (
	"sync"
	"os"
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"log"
	"io"
)

func NewFileQueue(prefix string, maxFileSize int64) *FileQueue {
	file := createFile(prefix, 0)
	fileWriter := bufio.NewWriterSize(file, 0)
	fr := createReader(file.Name())
	frb := bufio.NewReader(fr)

	fq := FileQueue{currentWF: file, currentW: fileWriter,
		currentRF: fr, currentR: frb,
		maxFileSize: maxFileSize, prefix: prefix}
	return &fq
}

type FileQueue struct {
	wMutex sync.Mutex
	rMutex sync.Mutex

	prefix      string
	currentNum  int64
	maxFileSize int64


	currentWF *os.File
	currentW  *bufio.Writer

	readersQueue chan string

	currentRF *os.File
	currentR  *bufio.Reader
}


func (f *FileQueue) Push(data []byte) {
	data = append(data, "\n"...)
	f.wMutex.Lock()
	defer f.wMutex.Unlock()
	f.currentW.Write(data)
	f.rotateLogFile()
}

//change the output to bytes.buffer channel in order to utilize less memory
func (f *FileQueue) Pop(n int64, s int64) io.Reader {
	f.rMutex.Lock()
	pr, pw := io.Pipe()

	go func(){
		defer f.rMutex.Unlock()
		for i := int64(0); i < n; i++ {
			line, err := f.currentR.ReadBytes('\n')
			pw.Write(line)
			if err != nil {
				fmt.Println(err)
				f.loadNewReader()
			}
		}
		pw.Close()
	}()

	return pr
}

func (f *FileQueue) Peek() (int64, int64){return 9999999,999999}
func (f *FileQueue) CanPush(s int, atomic bool) bool {return true}

func (f *FileQueue) Close(){
	f.rMutex.Lock()
	f.flushToDisk()
}


/**/
func (f *FileQueue) flushToDisk() {
	f.wMutex.Lock()
	defer f.wMutex.Unlock()
	f.currentW.Flush()
}

func (f *FileQueue) rotateLogFile(){
	fi, _ := f.currentWF.Stat()
	if fi.Size() >= f.maxFileSize {
		f.wMutex.Lock()
		defer f.wMutex.Unlock()
		f.currentNum++
		cw := f.currentWF.Name()
		f.currentWF.Close()
		go func() {
			f.readersQueue <- cw
		}()

		f.currentWF = createFile(f.prefix, f.currentNum)
		f.currentW = bufio.NewWriterSize(f.currentWF, 0)
	}
}

func (f *FileQueue) loadNewReader(){

	//delete file
	if f.currentWF.Name() != f.currentRF.Name() {
		f.currentRF.Close()
		os.Remove(f.currentRF.Name())
	}
	//pull new one from channel
	select {
	case x, ok := <-f.readersQueue:
		if ok {
			fmt.Println("Pulled new file")
			f.currentRF = createReader(x)
			f.currentR = bufio.NewReader(f.currentRF)
		}
	default:
		fmt.Println("No files to read from")
	}
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

func createReader(filename string) *os.File {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0777)
	if err != nil {
		log.Panic(err)
	}
	return f
}
