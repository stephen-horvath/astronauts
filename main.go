package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
)

type Astronaut struct {
	Name  string
	Craft string
}

// write astronauts to a csv file using the specified delimiter
func writeCsv(ast []Astronaut, delimiter string) {
	fo, err := os.Create("astronauts.csv")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	if _, err := fo.WriteString("Name" + delimiter + "Craft\n"); err != nil {
		panic(err)
	}
	for _, astronaut := range ast {
		if _, err := fo.WriteString(astronaut.Name + delimiter + astronaut.Craft + "\n"); err != nil {
			panic(err)
		}
	}
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
	const delimiter string = ","
	cbnPtr := flag.Bool("cbn", false, "Sort the Astronauts by Craft by Name")
	flag.Parse()
	cbn := *cbnPtr

	astroJson := fetch(url)
	astronauts := extractAstronauts(astroJson)
	sortAstronauts(astronauts, cbn)
	writeCsv(astronauts, delimiter)
}
