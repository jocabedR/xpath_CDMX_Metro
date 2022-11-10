package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/antchfx/xmlquery"
	"github.com/antchfx/xpath"
	"github.com/gorilla/mux"
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

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("{\"message\": \"Hello World!\"}")
}

func HandlePath(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Origin: %v\n", vars["origin"])
	fmt.Fprintf(w, "Destination: %v\n", vars["destination"])

	origin := vars["origin"]
	destination := vars["destination"]

	doc, root := readSource()

	readSource()

	countexpr, _ := xpath.Compile("count(//Folder[name='Estaciones de Metro']/Placemark)")
	numberStations := int(countexpr.Evaluate(xmlquery.CreateXPathNavigator(doc)).(float64))

	nodes, nodeNames := getNodes(root)
	segments := getSegments(root, nodes)

	_, okOrigin := nodes[origin]
	_, okDestination := nodes[destination]
	if !okOrigin || !okDestination {
		fmt.Fprintln(w, "Incorrect args.")
	} else {
		originID := nodes[origin].id
		destinationID := nodes[destination].id

		g := fillGraph(numberStations, segments, originID, destinationID, nodeNames)

		getPath(g, segments, nodeNames, originID, destinationID, w)
	}

}

func readSource() (*xmlquery.Node, xpath.NodeNavigator) {

	content, err := ioutil.ReadFile("Metro_CDMX.kml")
	if err != nil {
		log.Fatal(err)
	}

	xmlContent := string(content)

	doc, err := xmlquery.Parse(strings.NewReader(xmlContent))
	if err != nil {
		panic(err)
	}

	var root xpath.NodeNavigator
	root = xmlquery.CreateXPathNavigator(doc)

	return doc, root
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

func fillGraph(numberStations int, segments map[int]segment, o int, d int, nodeNames map[int]nodeName) *graph.Mutable {
	g := graph.New(numberStations)

	for i := 0; i < len(segments); i++ {
		g.AddBothCost(segments[i].origin, segments[i].destination, segments[i].distance)
	}

	return g
}

func getPath(g *graph.Mutable, segments map[int]segment, nodeNames map[int]nodeName, o int, d int, w http.ResponseWriter) {
	path, dist := graph.ShortestPath(g, o, d)

	if dist == -1 {
		fmt.Fprint(w, "Destination cannot be reached.\n")
	} else {
		fmt.Fprintln(w, "Number Stations: ", dist, " stations")
		fmt.Fprintln(w, "Line\t\tSegment")
		origin := path[0]
		for i := 1; i < len(path); i++ {
			destination := path[i]
			for j := 0; j < len(segments); j++ {
				if (segments[j].origin == origin && segments[j].destination == destination) || (segments[j].origin == destination && segments[j].destination == origin) {
					fmt.Fprintln(w, segments[j].line, "\t", nodeNames[origin].name, "-", nodeNames[destination].name)
				}
			}
			origin = destination
		}
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", HandleIndex).Methods(http.MethodGet)
	r.HandleFunc("/path/{origin}/{destination}", HandlePath)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Println("Listning at port " + srv.Addr)
	srv.ListenAndServe()
}
