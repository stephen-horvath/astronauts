package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Astronaut struct {
	Name  string
	Craft string
}

func extractAstronauts(astroJson []byte) []Astronaut {
	type Astronauts struct {
		People []Astronaut
	}
	var astronauts Astronauts
	err := json.Unmarshal(astroJson, &astronauts)
	if err != nil {
		panic(err)
	}
	return astronauts.People
}

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
	astroJson := fetch(url)
	astronauts := extractAstronauts(astroJson)
	fmt.Println(astronauts)
}
