package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/antchfx/xmlquery"
	"github.com/antchfx/xpath"
)

type node struct {
	id          int
	coordinates string
}

type segment struct {
	id       int
	origin   string
	destiny  string
	line     string
	distance float64
}

func main() {

	if len(os.Args) != 4 {
		fmt.Print("Incorrect number of arguments.\n")
		os.Exit(1)
	}

	fileName := os.Args[1]
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	xmlContent := string(content)

	origin := os.Args[2]
	destination := os.Args[3]

	doc, err := xmlquery.Parse(strings.NewReader(xmlContent))
	if err != nil {
		panic(err)
	}

	var root xpath.NodeNavigator
	root = xmlquery.CreateXPathNavigator(doc)

	nodes := getNodes(root)
	fmt.Println(nodes[origin].id)
	fmt.Println(nodes[destination].id)
	//segments := getSegments(root)
	//fmt.Println(segments)

}

func getNodes(root xpath.NodeNavigator) map[string]node {
	var nodes = make(map[string]node)

	expr := xpath.MustCompile("//Folder[name='Estaciones de Metro']/Placemark/name")
	stations := expr.Evaluate(root)
	iter := stations.(*xpath.NodeIterator)

	for counter := 0; iter.MoveNext(); counter++ {
		name := strings.Trim(iter.Current().Value(), " ")

		expr2 := xpath.MustCompile("//Folder[name='Estaciones de Metro']/Placemark[name='" + name + "']/Point/coordinates")
		coordinates := expr2.Evaluate(root)
		iter2 := coordinates.(*xpath.NodeIterator)

		for iter2.MoveNext() {
			coords := strings.Trim(strings.Trim(strings.Trim(iter2.Current().Value(), "\n"), " "), "\n")

			nodes[name] = node{
				counter, coords,
			}
		}
	}

	return nodes
}

/* func getSegments(root xpath.NodeNavigator) map[string]segment {
	var segments = make(map[string]segment)

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

		coor := expr3.Evaluate(root)
		var parseCor []string

		iter3 := coor.(*xpath.NodeIterator)
		for j := 1; iter3.MoveNext(); j++ {
			cordinates := strings.TrimRight(iter3.Current().Value(), " \n")
			parseCor = strings.Split(cordinates, "\n            ")
			parseCor = parseCor[1:]
		}

		for j := 0; j < len(parseCor); j++ {

			cordinate := parseCor[j]
			var nextcordinate string
			if j < len(parseCor)-1 {
				nextcordinate = parseCor[j+1]
			}

			expr2 := xpath.MustCompile("//Folder[name='Estaciones de Metro']/Placemark/Point[contains(coordinates,'" + cordinate + "')]/preceding-sibling::*[3]")
			stations := expr2.Evaluate(root)
			iter2 := stations.(*xpath.NodeIterator)

			for iter2.MoveNext() {
				station := iter2.Current().Value()
				fmt.Printf("\t%d. %s\n", j, station)
				fmt.Printf("\t|%s|\n", cordinate)

				var distanceNext float64
				if j < len(parseCor)-1 {
					distanceNext = distance(cordinate, nextcordinate)
					fmt.Printf("\t|%f|\n\n", distanceNext)

					segments[name] = node{
						/*id       int
						origin   string
						destiny  string
						line     string
						distance float64
						,
					}

				}

			}

		}

	}

	return segments
} */

func distance(coord0, coord1 string) float64 {
	coordenates0 := strings.Split(coord0, ",")
	coordenates1 := strings.Split(coord1, ",")

	x0, _ := strconv.ParseFloat(coordenates0[0], 64)
	x1, _ := strconv.ParseFloat(coordenates1[0], 64)
	y0, _ := strconv.ParseFloat(coordenates0[1], 64)
	y1, _ := strconv.ParseFloat(coordenates1[1], 64)

	return math.Sqrt(((x0 - x1) * (x0 - x1)) + ((y0 - y1) * (y0 - y1)))
}
