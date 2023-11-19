package main

import (
	"fmt"
	"io"
	"net/http"
)

// make an http GET request to the url returning the full response body as a slice
func fetch(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()
	if resp.StatusCode != 200 {
		panic(fmt.Sprintf("Houston we have a problem %s returned with status %d", url, resp.StatusCode))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return body
}

func main() {
	const url string = "http://api.open-notify.org/astros.json"
	json := fetch(url)
	fmt.Println(string(json))
}
