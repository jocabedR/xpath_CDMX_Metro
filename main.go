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
		fmt.Println("\tSTATIONS:")
		expr2 := xpath.MustCompile("//Folder[name='Estaciones de Metro']/Placemark[contains(description,'" + line + "')]/name")
		stations := expr2.Evaluate(root)
		iter2 := stations.(*xpath.NodeIterator)

		for i := 1; iter2.MoveNext(); i++ {
			station := iter2.Current().Value()
			fmt.Printf("\t%d. %s", i, strings.Trim(station, "\n"))

			expr3, err := xpath.Compile("//Folder[name='Estaciones de Metro']/Placemark[name='" + station + "']/Point/coordinates")
			if err != nil {
				panic(err)
			}

			var root2 xpath.NodeNavigator
			root2 = xmlquery.CreateXPathNavigator(doc)
			coor := expr3.Evaluate(root2)

			iter3 := coor.(*xpath.NodeIterator)
			for iter3.MoveNext() {
				fmt.Printf(", %s", strings.Trim(strings.Trim(iter3.Current().Value(), "\n"), " "))
			}
		}

	}

}
