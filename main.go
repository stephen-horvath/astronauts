package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
)

type Astronaut struct {
	Name  string
	Craft string
}

// sort the astronaughts. if cbn (craft by name) is true will be sorted by craft then name
// else sorted by name then craft
func sortAstronauts(ast []Astronaut, cbn bool) {
	sort.Slice(ast, func(i, j int) bool {
		if ast[i].Craft == ast[j].Craft && cbn || ast[i].Name != ast[j].Name && !cbn {
			return ast[i].Name < ast[j].Name
		} else {
			return ast[i].Craft < ast[j].Craft
		}
	})
}

// extracts a slice of Astronaut. The supplied byte slice is assumed to represent
// valid json with a key "people" containint he astronauts
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
	if resp.StatusCode != 201 {
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
	const cbn bool = false
	astroJson := fetch(url)
	astronauts := extractAstronauts(astroJson)
	sortAstronauts(astronauts, cbn)
	fmt.Println(astronauts)
}
