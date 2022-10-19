package main

import (
	"bufio"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func makeCall(url string, ch chan string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		ch <- "FAIL"
	}

	cl := http.Client{Timeout: 10 * time.Second}
	res, err := cl.Do(req)
	if err != nil {
		ch <- "FAIL"
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err == nil {
		ch <- "FAIL"
	}
	ch <- string(body)
}

func getQueries(ch chan []string) {
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

func main() {

}
