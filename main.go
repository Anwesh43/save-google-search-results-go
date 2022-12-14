package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func makeCall(url string, ch chan string, cb func()) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		ch <- "FAIL"
	}
	fmt.Println(url)

	cl := http.Client{Timeout: 10 * time.Second}
	res, err := cl.Do(req)
	if err != nil {
		ch <- "FAIL"
		cb()
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		ch <- "FAIL"
	}
	ch <- string(body)
	cb()
}

func getQueries(ch chan []string) {
	fmt.Println("Please start writing")
	words := make([]string, 0)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		word := scanner.Text()
		if word == "QUIT" {
			break
		}
		words = append(words, word)
	}
	ch <- words
}

func waitForAllCalls(urls []string, rangeCh chan string, ch chan string) {
	var wg sync.WaitGroup
	wg.Add(len(urls))
	for _, url := range urls {
		go makeCall(url, rangeCh, func() {
			wg.Done()
		})
	}
	ch <- "Done"
	wg.Wait()
	close(rangeCh)
}

func arrayMap(strArray []string, cb func(string) string) []string {
	newArr := make([]string, len(strArray))
	for _, word := range strArray {

		if strings.Trim(word, "") != "" {
			newArr = append(newArr, cb(word))
		}
	}
	return newArr
}

func makeQueriesCall(queries []string, writeCh chan string) {
	rangeCh := make(chan string, len(queries))
	chNext := make(chan string)
	results := make([]string, 0)
	go waitForAllCalls(arrayMap(queries, func(q string) string {
		return fmt.Sprintf("https://www.google.com/search?q=%s", q)
	}), rangeCh, chNext)
	<-chNext
	for currCh := range rangeCh {
		results = append(results, currCh)
		writeCh <- currCh
	}
	close(writeCh)
}

func createFile(ch chan string, doneChan chan string) {
	i := 0
	for result := range ch {
		if result == "FAIL" {
			continue
		}
		f, err := os.Create(fmt.Sprintf("search_%d.html", i))
		if err != nil {
			doneChan <- "FAIL"
		}
		f.Write(([]byte)(result))
		f.Close()
		i = i + 1
		fmt.Println("Result of ", i, "is written", (fmt.Sprintf("%d.txt", i)))
	}
	doneChan <- "Done Writing"
}

func main() {
	chArr := make(chan []string)

	go getQueries(chArr)
	words := <-chArr
	writeCh := make(chan string, len(words))
	doneCh := make(chan string)
	go makeQueriesCall(words, writeCh)
	go createFile(writeCh, doneCh)
	<-doneCh
}
