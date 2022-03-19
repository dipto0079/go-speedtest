package speedtest

import (
	"errors"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	chunkSize  = 16384
	windowSize = 10
)

var ErrTimeExpired = errors.New("time expired")


type Benchmark interface {
	Run(func(n int) error) error
}


type DownloadBenchmark struct {
	Client  http.Client
	Server  Server
	BaseURL string
}


func NewDownloadBenchmark(client http.Client, server Server) DownloadBenchmark {
	slashPos := strings.LastIndex(server.URL, "/")
	baseURL := server.URL[:slashPos] + "/random1000x1000.jpg"
	return DownloadBenchmark{client, server, baseURL}
}

func (b DownloadBenchmark) Run(fn func(n int) error) error {
	threadURL := b.BaseURL + "?x=" + strconv.Itoa(rand.Int())
	resp, err := b.Client.Get(threadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buf := make([]byte, chunkSize)
	for {
		num, err := resp.Body.Read(buf)
		nerr := fn(num)
		if nerr == ErrTimeExpired || err == io.EOF {
			break
		}
		if nerr != nil {
			return nerr
		}
		if err != nil {
			return err
		}
	}
	return nil
}


type UploadBenchmark struct {
	Client http.Client
	Server Server
}


func NewUploadBenchmark(client http.Client, server Server) UploadBenchmark {
	return UploadBenchmark{client, server}
}

func (b UploadBenchmark) Run(fn func(n int) error) error {
	reader := NewJunkReader(1024 * 1024)
	writer := NewCallbackWriter(fn)
	tee := io.TeeReader(&reader, writer)
	_, err := b.Client.Post(b.Server.URL, "text/plain", tee)
	return err
}

func RunBenchmark(b Benchmark, threads int, maxThreads int, duration time.Duration) int {
	var wg sync.WaitGroup
	var tc sync.Mutex

	
	resolution := time.Second / time.Duration(windowSize)
	chunks := make([]int, duration/resolution)

	
	reqs := make(chan int, maxThreads)
	for i := 0; i < threads; i++ {
		reqs <- 1
	}


	start := time.Now()
	active := true

	perform := func() {
		wg.Add(1)
		defer wg.Done()

		err := b.Run(func(n int) error {
			p := int(time.Since(start) / resolution)
			if p < len(chunks) {
				chunks[p] += n
			}
			if !active {
				return ErrTimeExpired
			}
			return nil
		})

		if active {
			if err != nil {
				log.Fatalln(err)
			}

			reqs <- 1

			tc.Lock()
			if threads < maxThreads {
				threads++
				reqs <- 1
			}
			tc.Unlock()
		} else {
			
			return
		}
	}

	
	timeout := time.After(duration)
	for active {
		select {
		case <-reqs:
			if active {
				go perform()
			}
		case <-timeout:
			
			active = false
		}
	}

	wg.Wait()

	maxSum := MaximalSumWindow(chunks, windowSize)
	windowAvg := MedianSumWindow(chunks, windowSize)
	return (maxSum + windowAvg) / 2
}
