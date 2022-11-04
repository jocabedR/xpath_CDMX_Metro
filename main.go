package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/antchfx/xmlquery"
	"github.com/antchfx/xpath"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Print("Incorrect number of arguments.\n")
		os.Exit(1)
	}

	fileName := os.Args[1]
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	s := string(content)

	doc, err := xmlquery.Parse(strings.NewReader(s))
	if err != nil {
		panic(err)
	}

	var root xpath.NodeNavigator
	root = xmlquery.CreateXPathNavigator(doc)

	fmt.Println("SUBWAY LINES:")
	expr := xpath.MustCompile("//Folder[name='Líneas de Metro']/Placemark/name")
	lines := expr.Evaluate(root)
	iter := lines.(*xpath.NodeIterator)

	for iter.MoveNext() {

		fmt.Println()
		line := iter.Current().Value()
		fmt.Println(line)

		expr3, err := xpath.Compile("//Folder[name='Líneas de Metro']/Placemark[name='" + line + "']/LineString/coordinates")
		if err != nil {
			panic(err)
		}

		var root2 xpath.NodeNavigator
		root2 = xmlquery.CreateXPathNavigator(doc)
		coor := expr3.Evaluate(root2)
		var parseCor []string

		iter3 := coor.(*xpath.NodeIterator)
		for i := 1; iter3.MoveNext(); i++ {
			cordinates := strings.TrimRight(iter3.Current().Value(), " \n")
			parseCor = strings.Split(cordinates, "\n            ")
			parseCor = parseCor[1:]
		}

		fmt.Println("\tSTATIONS:")
		for i := 0; i < len(parseCor); i++ {
			cordinate := parseCor[i]

			expr2 := xpath.MustCompile("//Folder[name='Estaciones de Metro']/Placemark/Point[contains(coordinates,'" + cordinate + "')]/preceding-sibling::*[3]")
			stations := expr2.Evaluate(root)
			iter2 := stations.(*xpath.NodeIterator)

			for iter2.MoveNext() {
				station := iter2.Current().Value()
				fmt.Printf("\t%d. %s, %s \n", i+1, station, cordinate)
			}

		}

	}

}
