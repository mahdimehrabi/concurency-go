package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

var websitesPool = make(chan Website, 10)

type Website struct {
	URL string
}

var random *rand.Rand

var urls = []string{"https://youtube.com",
	"https://facebook.com",
	"https://google.com", "https://stackoverflow.com", "https://github.com",
	"https://yahoo.com", "https://microsoft.com"}

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	random = rand.New(source)
}

func factory(ctx context.Context, num int, done chan bool) {
	defer func() {
		close(websitesPool)
		done <- true
	}()
	urlsLen := len(urls)
	for i := 0; i < num; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			websitesPool <- Website{
				urls[random.Intn(urlsLen)],
			}
		}
	}
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	start := time.Now().Unix()
	done := make(chan bool)
	wg := new(sync.WaitGroup)
	go factory(ctx, 50, done)
	go runWorker(wg, 3)

	go func() {
		for {
			fmt.Println("enter 'c' to stop:")
			var input string
			fmt.Scanln(&input)
			if input == "c" {
				cancel()
				fmt.Println("stopping...")
				return
			}
		}
	}()

	<-done
	fmt.Println("finished factory")
	wg.Wait()

	end := time.Now().Unix()
	fmt.Printf("it takes %d seconds\n", end-start)
}

func worker(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		website, open := <-websitesPool
		if !open && len(websitesPool) < 1 {
			return
		}
		downloadWrite(website)
	}
}

func runWorker(wg *sync.WaitGroup, num int) {
	wg.Add(num)
	for i := 0; i < num; i++ {
		go worker(wg)
	}
}

func downloadWrite(website Website) {
	time.Sleep(1 * time.Second)
	resp, err := http.Get(website.URL)
	if err != nil {
		fmt.Println("err:", err.Error())
		return
	}
	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, resp.Body)
	if err != nil {
		fmt.Println("err:", err.Error())
		return
	}
	f, err := os.OpenFile(fmt.Sprintf("websites/%d", rand.Intn(9999999999999999)), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("err:", err.Error())
		return
	}
	f.Write(buffer.Bytes())
}
