package main

import (
	"flag"
	"fmt"
	"github.com/pipiobjo/az-network-visualizer/config"
	"github.com/pipiobjo/az-network-visualizer/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func main() {
	InitLogger()

	var configFile string
	flag.StringVar(&configFile, "configFile", "", "Input file containing the azure network json")
	flag.Parse()

	required := []string{"configFile"}
	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	for _, req := range required {
		if !seen[req] {
			fmt.Fprintf(os.Stderr, "missing required -%s argument\n", req)
			os.Exit(2) // the same exit code flag.Parse uses
		}
	}
	MainController(configFile)

}

func InitLogger() {
	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Kitchen}
	log.Logger = zerolog.New(output).With().Caller().Timestamp().Logger()
}

func MainController(configFile string) {

	network := config.ReadInput(configFile)

	service.DotService(network)

}

// go run main.go | dot -Tpng  > test.png && open test.png

//func main() {
//	g := dot.NewGraph(dot.Directed)
//	n1 := g.Node("coding")
//	n2 := g.Node("testing a little").Box()
//
//	g.Edge(n1, n2)
//	g.Edge(n2, n1, "back").Attr("color", "red")
//
//	fmt.Println(g.String())
//}
