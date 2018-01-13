package adapters

import (
	"sync"
	"os"
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"log"
)

func NewFileQueue(prefix string, maxFileSize int64) *FileQueue {
	file := createFile(prefix, 0)
	fileWriter := bufio.NewWriterSize(file, 1)

	fq := FileQueue{currentWF: file, currentW: fileWriter,
		maxFileSize: maxFileSize, prefix: prefix, readersQueue: make(chan string, 99999)}
	go func(){fq.readersQueue <- file.Name()}()
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

	sizeBytes int64
	sizeCount int64
}

func (f *FileQueue) Push(data []byte) {
	data = append(data, "\n"...)
	f.wMutex.Lock()
	defer f.wMutex.Unlock()
	f.currentW.Write(data)
	f.currentW.Flush()
	f.sizeIncr(int64(len(data)))
	f.rotateLogFile()
}

//change the output to bytes.buffer channel in order to utilize less memory
func (f *FileQueue) Pop(out chan []byte, targetCount int64, targetSize int64) {

	for f.sizeCount > 0 && (targetCount > 0 || targetSize > 0) {
		if f.currentR == nil{
			if !f.loadNewReader() {
				break
			}
		}
		line, err := f.currentR.ReadBytes('\n')
		if err != nil {
			fmt.Println(err)
			if !f.loadNewReader() {
				break
			}

		} else {
			out <- line
			targetCount--
			s := int64(len(line))
			targetSize -= s
			f.sizeDecr(s)
		}

	}
	close(out)

}

func (f *FileQueue) Peek() (int64, int64) {
	f.rMutex.Lock()
	defer f.rMutex.Unlock()
	return f.sizeCount, f.sizeBytes
}
func (f *FileQueue) CanPush(s int64, atomic bool) bool { return true }

func (f *FileQueue) Close() {
	f.flushToDisk()
}

func (f *FileQueue) Prefix() string {
	return f.prefix
}

/**/
func (f *FileQueue) flushToDisk() {
	f.wMutex.Lock()
	defer f.wMutex.Unlock()
	f.currentW.Flush()
}

func (f *FileQueue) rotateLogFile() {
	fi, _ := f.currentWF.Stat()
	if fi.Size() >= f.maxFileSize {
		f.currentNum++
		cw := f.currentWF.Name()
		f.currentWF.Close()
		go func() {
			if f.currentNum != 1{ // dont send the first file twice
				f.readersQueue <- cw
			}
		}()

		f.currentWF = createFile(f.prefix, f.currentNum)
		f.currentW = bufio.NewWriterSize(f.currentWF, 1)
	}
}

func (f *FileQueue) loadNewReader() bool {
	r := false

	//delete file
	if f.currentRF != nil && f.currentWF.Name() != f.currentRF.Name() {
		f.currentRF.Close()
		os.Remove(f.currentRF.Name())
	}

	//pull new one from channel
	if len(f.readersQueue) > 0 {
		x := <- f.readersQueue
		fmt.Println("Pulled new file")
		f.currentRF = createReader(x)
		f.currentR = bufio.NewReader(f.currentRF)
		r = true
	} else {
		fmt.Println("No files to read from")
	}



	return r
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


func (f *FileQueue) sizeIncr(incr int64) {
	f.sizeBytes += incr
	f.sizeCount++
}

func (f *FileQueue) sizeDecr(decr int64) {
	f.sizeBytes -= decr
	f.sizeCount--
}
