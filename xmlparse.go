package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/twmb/algoimpl/go/graph"
)

type nd struct {
	XMLName xml.Name `xml:"nd"`
	Ref     string   `xml:"ref,attr"`
}

type tag struct {
	XMLName xml.Name `xml:"tag"`
	Kay     string   `xml:"k,attr"`
	Vee     string   `xml:"v,attr"`
}

type xmlway struct {
	XMLName xml.Name `xml:"way"`
	Refs    []nd     `xml:"nd"`
	Tags    []tag    `xml:"tag"`
}

type xmlnode struct {
	XMLName   xml.Name `xml:"node"`
	ID        string   `xml:"id,attr"`
	Latitude  float64  `xml:"lat,attr"`
	Longitude float64  `xml:"lon,attr"`
}

type osm struct {
	XMLName xml.Name  `xml:"osm"`
	Nodes   []xmlnode `xml:"node"`
	Ways    []xmlway  `xml:"way"`
}

func readXML(reader io.Reader) ([]xmlway, []xmlnode, error) {
	var osm osm
	if err := xml.NewDecoder(reader).Decode(&osm); err != nil {
		return nil, nil, err
	}

	return osm.Ways, osm.Nodes, nil
}

type idCount struct {
	ID    []string
	count []int
}

// We only want "ways" that have a named road otherwise it is an alley or private drive
func wayIsARoad(way xmlway) bool {
	var foundNamedRoad = false
	for i := range way.Tags {
		if way.Tags[i].Kay == "name" {
			foundNamedRoad = true
		}
	}

	return foundNamedRoad
}

func findIDAndAddToList(id string, idSlice *idCount) {
	var idFound = false

	for i := 0; i < len(idSlice.ID) && !idFound; i++ {
		if idSlice.ID[i] == id {
			idFound = true
			idSlice.count[i]++
		}
	}

	if !idFound {
		idSlice.ID = append(idSlice.ID, id)
		idSlice.count = append(idSlice.count, 1)
	}
}

// IDs (locations) with a count less than 2 are not at an intersection and we don't care about them
func removeIDThatHaveCountLessThanTwo(idSlice *idCount) {
	for i := range idSlice.ID {
		if idSlice.count[i] < 2 {
			idSlice.ID = append(idSlice.ID[:i], idSlice.ID[i+1:]...)
			idSlice.count = append(idSlice.count[:i], idSlice.count[i+1:]...)
		}
	}
}

func writeIntersectionsToFile(string fileName, idSlice idCount) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer file.Close()

	w := bufio.NewWriter(file)
	for i := range idSlice.ID {
		fmt.Fprintln(w, idSlice.ID[i])
	}

	return w.Flush()
}

func main() {
	nodesFile, err := filepath.Abs("neighborhood.xml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	file, err := os.Open(nodesFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	ways, nodes, err := readXML(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var listOfIDs idCount
	for i := range ways {
		for j := range ways[i].Refs {
			if wayIsARoad(ways[i]) {
				findIDAndAddToList(ways[i].Refs[j].Ref, &listOfIDs)
			}
		}
	}

	removeIDThatHaveCountLessThanTwo(&listOfIDs)
	writeIntersectionsToFile("listOfIntersections.txt", listOfIDs)

	count := 0
	for i := range listOfIDs.ID {
		if listOfIDs.count[i] > 1 {
			count++
			fmt.Printf("ID: %s Count: %d\n", listOfIDs.ID[i], listOfIDs.count[i])
		}
	}

	fmt.Printf("Number of intersections: %d\n", count)
	fmt.Printf("Number of nodes: %d\n", len(nodes))

	g := graph.New(graph.Undirected)
	nodes := make(map[rune]graph.Node, 0)
}
