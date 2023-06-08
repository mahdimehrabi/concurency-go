package main

import (
	"bytes"
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

func factory(num int, done chan bool) {
	urlsLen := len(urls)
	for i := 0; i < num; i++ {
		websitesPool <- Website{
			urls[random.Intn(urlsLen)],
		}
	}
	close(websitesPool)
	done <- true
}

func main() {
	start := time.Now().Unix()
	done := make(chan bool)
	go factory(50, done)
	runWorker(3)
	<-done

	end := time.Now().Unix()
	fmt.Printf("it takes %d seconds\n", end-start)
}

func worker(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		website, open := <-websitesPool
		if !open {
			return
		}
		downloadWrite(website)
	}
}

func runWorker(num int) {
	wg := new(sync.WaitGroup)
	wg.Add(num)
	for i := 0; i < num; i++ {
		go worker(wg)
	}
	wg.Wait()
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
