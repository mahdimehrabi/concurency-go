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

func factory(num int) []Website {
	websites := make([]Website, num)
	urlsLen := len(urls)
	for i := 0; i < num; i++ {
		websites[i] = Website{
			urls[random.Intn(urlsLen)],
		}
	}
	return websites
}

func main() {
	start := time.Now().Unix()
	wg := new(sync.WaitGroup)
	for _, website := range factory(25) {
		wg.Add(1)
		go downloadWrite(website, wg)
	}
	wg.Wait()
	end := time.Now().Unix()
	fmt.Printf("it takes %d seconds\n", end-start)
}

func downloadWrite(website Website, wg *sync.WaitGroup) {
	defer wg.Done()
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
	defer f.Close()
	f.Write(buffer.Bytes())
}
