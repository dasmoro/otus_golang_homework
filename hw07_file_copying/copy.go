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

	if !fi.Mode().IsRegular() {
		return errors.New("unsupported type")
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
					DrawGopher(x, y+2, i, percents)
					time.Sleep(100 * time.Millisecond)
				}
				time.Sleep(200 * time.Millisecond)
				up = false
			} else {
				for i := 10; i >= 0; i-- {
					DrawGopher(x, y+2, i, percents)
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

func DrawGopher(x int, y int, frame int, percents *byte) {
	termMu.Lock()
	log.Printf("\033[%d;%dH", y, x) // place cursor at top left corner
	log.Printf("%v%% complete\n", *percents)

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

	termMu.Unlock()
}

func readPipe(pw io.PipeWriter, fromFile *MyFile, needToRead *int64) {
	defer pw.Close()
	buf := make([]byte, chunkSize)
	localNeedToRead := atomic.LoadInt64(needToRead)
	for {
		if localNeedToRead < int64(chunkSize) {
			buf = buf[0:int(localNeedToRead)]
		}
		n, err := fromFile.file.Read(buf)
		localNeedToRead -= int64(n)
		atomic.StoreInt64(needToRead, localNeedToRead)
		_, pErr := pw.Write(buf)

		if pErr != nil {
			break
		}
		if errors.Is(err, io.EOF) || localNeedToRead <= 0 {
			break
		}
	}
}

func progress(fromSize int64, elapsed *int64) {
	size := float64(fromSize)
	var percents byte
	doneCh := make(chan interface{})
	go RenderProgress(0, 0, &percents, doneCh)
	for {
		time.Sleep(time.Millisecond * 100)
		n := float64(atomic.LoadInt64(elapsed))
		percents = 100 - byte(n/size*100)
		if percents >= 100 {
			doneCh <- 1
			break
		}
	}
}

func prepare(fromPath, toPath string, offset int64) (*MyFile, *MyFile, error) {
	if fromPath == toPath {
		return nil, nil, errors.New("from and to are equals")
	}
	from, err := NewMyFile(fromPath, os.O_RDONLY)
	if err != nil {
		return nil, nil, err
	}
	err = from.SetOffset(offset)
	if err != nil {
		return nil, nil, err
	}
	to, err := NewMyFile(toPath, os.O_CREATE|os.O_RDWR)
	if err != nil {
		return nil, nil, err
	}
	err = to.file.Truncate(0)
	if err != nil {
		return nil, nil, err
	}
	return from, to, nil
}

func Copy(fromPath, toPath string, offset, limit int64, showProgress bool) error {
	from, to, err := prepare(fromPath, toPath, offset)
	if err != nil {
		return err
	}

	fromSize := from.GetSize()
	needToRead := fromSize - offset
	if limit < needToRead && limit != 0 {
		needToRead = limit
	}

	defer func() {
		from.file.Close()
		to.file.Close()
	}()

	pr, pw := io.Pipe()

	go readPipe(*pw, from, &needToRead)

	if showProgress {
		go progress(fromSize, &needToRead)
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
		log.Print("\n100% complete\n")
	}

	return nil
}
