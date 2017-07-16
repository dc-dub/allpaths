package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type ND struct {
	XMLName xml.Name `xml:"nd"`
	Ref     string   `xml:"ref,attr"`
}

type Tag struct {
	XMLName xml.Name `xml:"tag"`
	Kay     string   `xml:"k,attr"`
	Vee     string   `xml:"v,attr"`
}

type XMLWay struct {
	XMLName xml.Name `xml:"way"`
	Refs    []ND     `xml:"nd"`
	Tags    []Tag    `xml:"tag"`
}

type XMLNode struct {
	XMLName   xml.Name `xml:"node"`
	ID        string   `xml:"id,attr"`
	Latitude  float64  `xml:"lat,attr"`
	Longitude float64  `xml:"lon,attr"`
}

type OSM struct {
	XMLName xml.Name  `xml:"osm"`
	Nodes   []XMLNode `xml:"node"`
	Ways    []XMLWay  `xml:"way"`
}

func ReadXML(reader io.Reader) ([]XMLWay, []XMLNode, error) {
	var osm OSM
	if err := xml.NewDecoder(reader).Decode(&osm); err != nil {
		return nil, nil, err
	}

	return osm.Ways, osm.Nodes, nil
}

type NodeSlice []XMLNode

func (slice NodeSlice) Len() int           { return len(slice) }
func (slice NodeSlice) Less(i, j int) bool { return slice[i].ID > slice[j].ID }
func (slice NodeSlice) Swap(i, j int)      { slice[i], slice[j] = slice[j], slice[i] }

type IDCount struct {
	ID    []string
	count []int
	Roads []string
}

func FindID(id string, idSlice *IDCount) {
	var idFound bool = false

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

	xmlWay, xmlNodes, err := ReadXML(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var listOfIDs IDCount
	for i := range xmlWay {
		for j := range xmlWay[i].Refs {
			FindID(xmlWay[i].Refs[j].Ref, &listOfIDs)
		}
	}

	count := 0
	for i := range listOfIDs.ID {
		if listOfIDs.count[i] > 1 {
			count++
			fmt.Printf("ID: %s Count: %d\n", listOfIDs.ID[i], listOfIDs.count[i])
		}
	}
	fmt.Printf("Number of intersections: %d\n", count)
	fmt.Printf("Number of nodes: %d\n", len(xmlNodes))
}
