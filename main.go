package main

import (
	"io/ioutil"
	"net/http"
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

func main() {

}
