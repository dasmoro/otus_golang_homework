package main

import (
	"errors"
	"io"
	"log"
	"os"
	"strings"
	"sync"
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

func (ef *MyFile) Init(filename string, mode int) error {
	file, err := os.OpenFile(filename, mode, 0644)
	if err != nil {
		return err
	}

	fi, err := file.Stat()
	if err != nil {
		file.Close()
		return err
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

func BuildMyFile(filename string, mode int) (*MyFile, error) {
	f := MyFile{}
	err := f.Init(filename, mode)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

var termMu = sync.Mutex{}

func RenderProgress(x int, y int, percents *byte, done <-chan interface{}) {
	log.Print("\033[2J") // clear terminal
	up := true
	for {
		select {
		case <-done:
			break
		default:
			if up {
				for i := 0; i < 11; i++ {
					termMu.Lock()
					log.Printf("\033[%d;%dH", y, x) // place cursor at top left corner
					log.Printf("%v%% complete\n", *percents)
					DrawGopher(x, y+2, i)
					termMu.Unlock()
					time.Sleep(100 * time.Millisecond)
				}
				time.Sleep(200 * time.Millisecond)
				up = false
			} else {
				for i := 10; i >= 0; i-- {
					termMu.Lock()
					log.Printf("\033[%d;%dH", y, x) // place cursor at top left corner
					log.Printf("%v%% complete\n", *percents)
					DrawGopher(x, y+2, i)
					termMu.Unlock()
					time.Sleep(30 * time.Millisecond)
				}
				time.Sleep(200 * time.Millisecond)
				log.Print("\033[2J")
				x += 5
				up = true
			}
		}
	}
}

func DrawGopher(x int, y int, frame int) {
	gopher := []string{
		`        ,_---~~~~~----._         `,
		`  _,,_,*^____      _____''*g*\\\"*, `,
		` / __/ /'     ^.  /      \ ^@q   f `,
		`[  @f | @))    |  | @))   l  0 _/  `,
		` \'/   \~____ / __ \_____/    \   `,
		`  |           _l__l_           I  `,
		`  }          [______]           I  `,
		`  ]            | | |            |  `,
		`  ]             ~ ~             |  `,
		`  |                             |   `,
		`   |                           |   `,
	}
	emptyStr := `                                   `
	log.Printf("\033[%d;%dH", y, 0)
	for i := 0; i < len(gopher); i++ {
		if i < len(gopher)-1-frame {
			log.Printf("%v%v\n", strings.Repeat(" ", x), emptyStr)
		} else {
			log.Printf("%v%v\n", strings.Repeat(" ", x), gopher[i-len(gopher)+1+frame])
		}
	}
}

func Copy(fromPath, toPath string, offset, limit int64, showProgress bool) error {
	from, err := BuildMyFile(fromPath, os.O_RDONLY)

	if err != nil {
		return err
	}

	err = from.SetOffset(offset)
	if err != nil {
		return err
	}
	to, err := BuildMyFile(toPath, os.O_CREATE|os.O_RDWR)
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

	go func() {
		defer pw.Close()
		buf := make([]byte, chunkSize)
		for {
			if N < int64(chunkSize) {
				buf = buf[0:N]
			}
			n, err := from.file.Read(buf)

			N -= int64(n)
			_, pErr := pw.Write(buf)

			if pErr != nil {
				break
			}

			if errors.Is(err, io.EOF) || N <= 0 {
				break
			}
		}
	}()

	if showProgress {
		go func() {
			size := float64(atomic.LoadInt64(&fromSize))
			var percents byte = 0
			doneCh := make(chan interface{})
			go RenderProgress(0, 0, &percents, doneCh)
			for {
				time.Sleep(time.Millisecond * 100)
				n := float64(atomic.LoadInt64(&N))
				percents = 100 - byte(n/size*100)
				if percents >= 100 {
					doneCh <- 1
					break
				}
			}
		}()
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
		log.Print("\033[2J")
		log.Printf("\033[%d;%dH", 0, 0)
		log.Print("100% complete")
	}

	return nil
}
