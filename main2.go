package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/antchfx/xmlquery"
	"github.com/antchfx/xpath"
	"github.com/yourbasic/graph"
)

type node struct {
	id          int
	coordinates string
}

type nodeName struct {
	name string
}

type segment struct {
	name        string
	line        string
	origin      int
	destination int
	distance    int64
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
	countexpr, err := xpath.Compile("count(//Folder[name='Estaciones de Metro']/Placemark)")
	numberStations := int(countexpr.Evaluate(xmlquery.CreateXPathNavigator(doc)).(float64))

	var root xpath.NodeNavigator
	root = xmlquery.CreateXPathNavigator(doc)

	nodes, nodeNames := getNodes(root)
	segments := getSegments(root, nodes)

	fillGraph(numberStations, segments, nodes[origin].id, nodes[destination].id, nodeNames)
}

func getNodes(root xpath.NodeNavigator) (map[string]node, map[int]nodeName) {
	var nodes = make(map[string]node)
	var nodeNames = make(map[int]nodeName)

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

			nodes[name] = node{counter, coords}

			nodeNames[counter] = nodeName{name}
		}
	}
	return nodes, nodeNames
}

func getSegments(root xpath.NodeNavigator, nodes map[string]node) map[int]segment {
	var segments = make(map[int]segment)

	// GET LINES
	expr := xpath.MustCompile("//Folder[name='Líneas de Metro']/Placemark/name")
	lines := expr.Evaluate(root)
	iterLines := lines.(*xpath.NodeIterator)

	counter := 0

	for iterLines.MoveNext() {

		line := iterLines.Current().Value()
		// GET STATIONS COORDINATES FOR THE CURRENT LINE
		expr3 := xpath.MustCompile("//Folder[name='Líneas de Metro']/Placemark[name='" + line + "']/LineString/coordinates")

		coor := expr3.Evaluate(root)
		var parseCor []string

		iter3 := coor.(*xpath.NodeIterator)
		for j := 1; iter3.MoveNext(); j++ {
			cordinates := strings.TrimRight(iter3.Current().Value(), " \n")
			parseCor = strings.Split(cordinates, "\n            ")
			parseCor = parseCor[1:]
		}

		//fmt.Println()
		for j := 0; j < len(parseCor); j++ {

			cordinate := parseCor[j]
			var nextcordinate = ""

			if j < len(parseCor)-1 {
				nextcordinate = parseCor[j+1]
			}

			// GET STATION NAME FOR THE CURRENT COORDINATES
			expr2 := xpath.MustCompile("//Folder[name='Estaciones de Metro']/Placemark/Point[contains(coordinates,'" + cordinate + "')]/preceding-sibling::*[3]")
			stations := expr2.Evaluate(root)
			iter2 := stations.(*xpath.NodeIterator)

			for iter2.MoveNext() {
				station := iter2.Current().Value()

				if j < len(parseCor)-1 {
					var destination string

					// GET STATION NAME FOR THE NEXT COORDINATES
					expr3 := xpath.MustCompile("//Folder[name='Estaciones de Metro']/Placemark/Point[contains(coordinates,'" + nextcordinate + "')]/preceding-sibling::*[3]")
					dest := expr3.Evaluate(root)
					iter3 := dest.(*xpath.NodeIterator)

					for iter3.MoveNext() {
						destination = iter3.Current().Value()
					}

					if destination != "" {
						name := station + "-" + destination + " / " + destination + station
						segments[counter] = segment{
							name, line, nodes[station].id, nodes[destination].id, 1,
						}
						counter++
					}

				}
			}
		}
	}
	return segments
}

func fillGraph(numberStations int, segments map[int]segment, o int, d int, nodeNames map[int]nodeName) {
	g := graph.New(numberStations)

	for i := 0; i < len(segments); i++ {
		g.AddBothCost(segments[i].origin, segments[i].destination, segments[i].distance)
	}

	path, dist := graph.ShortestPath(g, o, d)

	if dist == -1 {
		fmt.Print("Destination cannot be reached.\n")
		os.Exit(1)
	}

	fmt.Println("Number Stations: ", dist, " stations")
	getStations(path, segments, nodeNames)
}

func getStations(path []int, segments map[int]segment, nodeNames map[int]nodeName) {
	origin := path[0]
	for i := 1; i < len(path); i++ {
		destination := path[i]
		for j := 0; j < len(segments); j++ {
			if (segments[j].origin == origin && segments[j].destination == destination) || (segments[j].origin == destination && segments[j].destination == origin) {
				fmt.Println(segments[j].line, " ", nodeNames[origin].name, "-", nodeNames[destination].name)
			}
		}
		origin = destination
	}
}
