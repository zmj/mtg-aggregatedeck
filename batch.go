package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func splitpath(path string) (string, string) {
	path_components := strings.Split(path, "/")
	basename := path_components[len(path_components)-1]
	dirname := strings.Join(path_components[:len(path_components)-1], "/")
	return basename, dirname
}

func is_hidden_file(path string) bool {
	basename, _ := splitpath(path)
	return strings.HasPrefix(basename, ".")
}

func aggregate_to_file(decks []*Deck, deckname string, output_path string) {
	var archetypical_deck *Deck
	if len(decks) == 1 {
		// trivially, a single deck's aggregate is itself
		archetypical_deck = decks[0]
	} else {
		var err error
		archetypical_deck, err = aggregate(decks)
		if err != nil {
			log.Fatalln("Couldn't create archetype for", deckname+":", err)
		}
	}

	err := os.MkdirAll(output_path, 0755)
	if err != nil {
		log.Fatalln("Couldn't create output directory", output_path+":", err)
	}
	filename := output_path + "/" + deckname + "-aggregate.txt"
	err = ioutil.WriteFile(filename, []byte(archetypical_deck.String()), 0644)
	if err != nil {
		log.Fatalln("Couldn't write output file", output_path+":", err)
	}
}

func strip_trailing_slash(s *string) {
	l := len(*s)
	if (*s)[l-1] == '/' {
		*s = (*s)[:l-2]
	}
}

func batch(top_path string, output_path string, verbose bool) {
	decks_processed := 0
	aggregate_decks_created := 0

	strip_trailing_slash(&top_path)
	strip_trailing_slash(&output_path)

	var deck_dir, deckname string
	decks := make([]*Deck, 0)

	filepath.Walk(top_path, func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			return nil
		}

		if is_hidden_file(path) {
			if verbose {
				log.Println("  Ignoring hidden file", path)
			}
			return nil
		}

		basename, dirname := splitpath(path)

		// if we're traversing a new folder
		if deck_dir != dirname {

			// aggregate what we've got already from previous folder
			if len(decks) > 0 {
				deckname := strings.Replace(deck_dir[len(top_path)+1:], "/", "-", -1)

				if verbose {
					log.Println("  Aggregating previous", len(decks), "decks")
				}
				aggregate_to_file(decks, deckname, output_path)
				aggregate_decks_created++
			}

			// allocate structures to collect decks from this new folder
			decks = make([]*Deck, 0)
			deck_dir = dirname
			deckname = strings.Replace(deck_dir[len(top_path)+1:], "/", "-", -1)
			if verbose {
				log.Println("Scanning for decks in", dirname)
			}
		}

		// we've found a deckfile
		if verbose {
			log.Println("  Parsing deck", basename)
		}

		// read the deckfile
		deckfile, err := os.Open(path)
		if err != nil {
			log.Fatalln("Path open failed for deck", path+":", err)
		}
		deck, err := NewDeck(deckfile)
		if err != nil {
			log.Fatalln("Deck parsing failed for deck", path+":", err)
		}

		decks = append(decks, deck)
		decks_processed++
		return nil
	})

	if verbose {
		log.Println("  Aggregating previous", len(decks), "decks")
	}
	aggregate_to_file(decks, deckname, output_path)
	aggregate_decks_created++

	log.Println(decks_processed, "decks processed from", top_path)
	log.Println(aggregate_decks_created, "aggregate decks created in", output_path)
}
