package main

import (
	"flag"
	"log"
	"strconv"
)

var batchmode = flag.Bool("batch", false, "Whether to run in batch mode (requires -decks and -out)")
var deckpath = flag.String("decks", "", "(batch mode) A directory with subdirectories containing deck files")
var outpath = flag.String("out", "", "(batch mode) Output directory for aggregated decks; will be created if it doesn't exist")
var servermode = flag.Bool("web", false, "Whether to run in webserver mode")
var port = flag.Int("port", 8981, "(web mode) Port to bind to")
var prefix = flag.String("prefix", "metadeck", "(web mode) Route at which to serve the web interface")
var verbose = flag.Bool("v", false, "Give more verbose output")

func main() {
	flag.Parse()
	if *batchmode || (*deckpath != "" && *outpath != "") {
		if *deckpath == "" || *outpath == "" {
			log.Fatalln("Batch operation requires --decks and --out arguments. Consult --help for more information.")
		}
		batch(*deckpath, *outpath, *verbose)
	} else {
		// servermode
		bind := "127.0.0.1:" + strconv.Itoa(*port)
		log.Println("Starting webserver on", bind+"/"+*prefix)
		run_server(bind, *prefix)
	}
}
