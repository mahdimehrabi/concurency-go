package main

import (
	"context"
	"fmt"
	"time"
)

func download(url string, done chan bool) {
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		fmt.Println("downloading", url)
	}
	done <- true
}

func RunDownload(ctx context.Context, done chan bool) {
	downloadDone := make(chan bool)
	go download("google.com", downloadDone)
	select {
	case <-ctx.Done():
		fmt.Println("canceled")
	case <-downloadDone:
		fmt.Println("download completed")
	}
	done <- true
}

func main() {
	ctx := context.Background()
	rndDone := make(chan bool)
	ctx, _ = context.WithTimeout(ctx, 11*time.Second)
	go RunDownload(ctx, rndDone)

	<-rndDone
}
