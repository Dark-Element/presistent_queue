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
	file := createFile(prefix, 0)
	fileWriter := bufio.NewWriterSize(file, 1024*1024*5)

	fr := createReader(file.Name())
	frb := bufio.NewReader(fr)

	fq := FileQueue{currentWF: file, currentW: fileWriter,
		currentRF: fr, currentR: frb,
		maxFileSize: maxFileSize, prefix: prefix}
	return &fq
}

type iFile interface {
	Push(data bytes.Buffer, flush bool)
	Pop(c int64) bytes.Buffer
	Close()
}

type FileQueue struct {
	wMutex sync.Mutex
	rMutex sync.Mutex

	currentWF *os.File
	currentW  *bufio.Writer

	readersQueue chan string

	currentRF *os.File
	currentR  *bufio.Reader

	prefix      string
	currentNum  int64
	maxFileSize int64
}

//push to file
//lock the request until flushed (ACID compliance)
func (f *FileQueue) Push(data []byte, flush bool) {
	data = append(data, "\n"...)
	f.wMutex.Lock()
	f.currentW.Write(data)
	f.wMutex.Unlock()
	if flush {
		f.flushToDisk()
	}
}

func (f *FileQueue) flushToDisk() {
	f.wMutex.Lock()
	//f.currentW.Flush()
	fi, err := f.currentWF.Stat()
	if err != nil {
		log.Panic("BBBBBBBBBBB")
	}
	if fi.Size() >= f.maxFileSize {
		f.currentNum++
		cw := f.currentWF.Name()
		f.currentWF.Close()
		go func() {
			f.readersQueue <- cw
		}()

		f.currentWF = createFile(f.prefix, f.currentNum)
		f.currentW = bufio.NewWriterSize(f.currentWF, 1024*1024*5)
	}
	f.wMutex.Unlock()
}

//change the output to bytes.buffer channel in order to utilize less memory
func (f *FileQueue) Pop(c int) bytes.Buffer {
	f.rMutex.Lock()
	defer f.rMutex.Unlock()
	b := bytes.Buffer{}
	var line string
	var err error
	for i := 0; i < c; i++ {
		line, err = f.currentR.ReadString('\n')
		b.WriteString(line)
		if err != nil {
			fmt.Println(err)
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
				return b
			}
		}
	}

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

func createReader(filename string) *os.File {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0777)
	if err != nil {
		log.Panic(err)
	}
	return f
}
