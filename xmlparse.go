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

type OSM struct {
	XMLName xml.Name `xml:"osm"`
	Ways    []XMLWay `xml:"way"`
}

func ReadXML(reader io.Reader) ([]XMLWay, error) {
	var osm OSM
	if err := xml.NewDecoder(reader).Decode(&osm); err != nil {
		return nil, err
	}

	return osm.Ways, nil
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

	xmlWay, err := ReadXML(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(len(xmlWay))
	num := 0
	for i := 0; i < len(xmlWay); i++ {
		for j := 0; j < len(xmlWay[i].Tags); j++ {
			streetName := xmlWay[i].Tags[j].Vee
			if xmlWay[i].Tags[j].Kay == "name" {
				fmt.Printf("Street Name: %s\n", streetName)
			}
			if streetName == "Winterlake Drive" {
				fmt.Println("We're HOME!")

			}
		}
	}
	fmt.Printf("\n\n Total number of marshes: %d", num)
}
