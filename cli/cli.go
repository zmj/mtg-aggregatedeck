package main

import (
	"flag"
	"log"
)

var deckpath = flag.String("decks", "", "(batch mode) A directory with subdirectories containing deck files")
var outpath = flag.String("out", "", "(batch mode) Output directory for aggregated decks; will be created if it doesn't exist")
var verbose = flag.Bool("v", false, "Give more verbose output")

func main() {
	flag.Parse()
	if *deckpath == "" || *outpath == "" {
		log.Fatalln("Batch operation requires --decks and --out arguments. Consult --help for more information.")
	}
	batch(*deckpath, *outpath, *verbose)
}
