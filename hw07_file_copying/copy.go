package main

import (
	"errors"
	"io"
	"log"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

var (
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	chunkSize                = 50 // more is faster, but I want see a gopher :)
)

type MyFile struct {
	filename string
	file     *os.File
	fi       os.FileInfo
	size     int64
}

func (ef *MyFile) init(filename string, mode int) error {
	file, err := os.OpenFile(filename, mode, 0644)
	if err != nil {
		return err
	}

	fi, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	if fi.IsDir() {
		return errors.New("unsupported type (directory)")
	}
	ef.filename = filename
	ef.file = file
	ef.fi = fi
	return nil
}

func (ef *MyFile) GetSize() int64 {
	if ef.size != 0 {
		return ef.size
	}
	if ef.fi != nil {
		ef.size = ef.fi.Size()
		return ef.size
	}
	return 0
}

func (ef *MyFile) SetOffset(offset int64) error {
	if offset > ef.GetSize() {
		ef.file.Close()
		return ErrOffsetExceedsFileSize
	}
	_, err := ef.file.Seek(offset, io.SeekStart)
	if err != nil {
		return nil
	}
	return nil
}

func NewMyFile(filename string, mode int) (*MyFile, error) {
	f := MyFile{}
	err := f.init(filename, mode)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func readPipe(pw io.PipeWriter, fromFile *MyFile, needToRead *int64) {
	defer pw.Close()
	buf := make([]byte, chunkSize)
	for {
		if *needToRead < int64(chunkSize) {
			buf = buf[0:*needToRead]
		}
		n, err := fromFile.file.Read(buf)

		*needToRead -= int64(n)
		_, pErr := pw.Write(buf)

		if pErr != nil {
			break
		}

		if errors.Is(err, io.EOF) || *needToRead <= 0 {
			break
		}
	}
}

func progress(fromSize int64, elapsed *int64) {
	log.Print("\033[2J") // clear terminal
	log.Printf("\033[%d;%dH", 2, 1)
	log.Print("[                    ]")
	log.Printf("\033[%d;%dH", 1, 1)
	size := float64(atomic.LoadInt64(&fromSize))
	var percents byte
	var barTick byte
	for {
		log.Printf("%d%% complete", percents)
		if percents > 0 {
			barTick = percents / 5
			bar := "[" + strings.Repeat("=", int(barTick)) + ">"
			log.Printf("\033[%d;%dH", 2, 1)
			log.Print(bar)
		}
		log.Printf("\033[%d;%dH", 1, 1)
		time.Sleep(time.Millisecond * 100)
		n := float64(atomic.LoadInt64(elapsed))
		percents = 100 - byte(n/size*100)
		if percents >= 100 {
			break
		}
	}
}

func Copy(fromPath, toPath string, offset, limit int64, showProgress bool) error {
	if fromPath == toPath {
		return errors.New("from and to are equals")
	}
	from, err := NewMyFile(fromPath, os.O_RDONLY)

	if err != nil {
		return err
	}

	err = from.SetOffset(offset)
	if err != nil {
		return err
	}
	to, err := NewMyFile(toPath, os.O_CREATE|os.O_RDWR)
	if err != nil {
		return err
	}

	err = to.file.Truncate(0)
	if err != nil {
		return err
	}

	fromSize := from.GetSize()
	N := fromSize - offset
	if limit < N && limit != 0 {
		N = limit
	}

	defer func() {
		from.file.Close()
		to.file.Close()
	}()

	pr, pw := io.Pipe()

	go readPipe(*pw, from, &N)

	if showProgress {
		go progress(fromSize, &N)
	}

	buf := make([]byte, chunkSize)
	for {
		n, err := pr.Read(buf)
		_, wErr := to.file.Write(buf[:n])
		if wErr != nil {
			break
		}
		if errors.Is(err, io.EOF) {
			break
		}
	}
	pr.Close()

	if showProgress {
		log.Print("\n\n100% complete\n")
	}

	return nil
}
