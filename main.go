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

	fmt.Println("Lines:")
	expr := xpath.MustCompile("//Folder[name='Líneas de Metro']/Placemark/name")
	lines := expr.Evaluate(root)
	iter := lines.(*xpath.NodeIterator)

	for iter.MoveNext() {

		fmt.Println()
		line := iter.Current().Value()
		fmt.Println(line)
		fmt.Println("\tStations:")
		expr2 := xpath.MustCompile("//Folder[name='Estaciones de Metro']/Placemark[contains(description,'" + line + "')]/name")
		stations := expr2.Evaluate(root)
		iter2 := stations.(*xpath.NodeIterator)

		for iter2.MoveNext() {
			station := iter2.Current().Value()
			fmt.Println("\t" + station)
			/* expr = xpath.MustCompile("//Folder[name='Estaciones de Metro']/Placemark[name=" + station + "]/Point/coordinates")
			val := expr.Evaluate(root)
			iter := val.(*xpath.NodeIterator)
			for iter.MoveNext() {
				fmt.Print(iter.Current().Value())
			} */

			/* expr, err := xpath.Compile("//Folder[name='Estaciones de Metro']/Placemark[name=" + station + "]/Point/coordinates")
			if err != nil {
				panic(err)
			}
			val := expr.Evaluate(root)
			fmt.Print(val) */
		}

	}

}
