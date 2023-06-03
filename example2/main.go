package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

type Post struct {
	UserId int64  `json:"user_id"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func main() {
	msg := make(chan Post)
	wg := &sync.WaitGroup{}

	wg.Add(1)

	go sendMessage(msg, wg)
	readMessage(msg)

	wg.Wait()
}

func sendMessage(ch chan<- Post, wg *sync.WaitGroup) {

	var post Post

	for i := 0; i < 100; i++ {

		res, err := http.Get("https://jsonplaceholder.typicode.com/posts/" + fmt.Sprint(i))

		if err != nil {
			log.Fatal(err)
		}

		defer res.Body.Close()

		bt, err := io.ReadAll(res.Body)

		if err != nil {
			log.Fatal(err)
		}

		if err = json.Unmarshal(bt, &post); err != nil {
			log.Fatal(err)
		}

		ch <- post

	}
	close(ch)

	defer wg.Done()
}

func readMessage(ch <-chan Post) {
	for v := range ch {
		fmt.Println(v)
	}
}
