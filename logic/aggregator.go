package main

import "sort"
import "errors"
import "fmt"

var maindecksize = 60
var sideboardsize = 15

func aggregate(decks []*Deck) (*Deck, error) {
	maindecks := make([][]*Card, len(decks))
	sideboards := make([][]*Card, len(decks))
	for i, deck := range decks {
		maindecks[i] = deck.maindeck
		sideboards[i] = deck.sideboard
	}
	maindeck, err := agg(maindecks, maindecksize)
	if err != nil {
		return nil, err
	}
	sideboard, err := agg(sideboards, sideboardsize)
	if err != nil {
		return nil, err
	}
	return &Deck{maindeck: maindeck, sideboard: sideboard}, nil
}

type Metacard struct {
	name     string
	instance int // nth copy of card
	total    int // total # of card played in all decklists
	count    int // # of nth copy in all decklists
}

func agg(cardlists [][]*Card, decksize int) ([]*Card, error) {
	metalist := map[Card]*Metacard{}
	for _, list := range cardlists {
		for _, card := range list {
			for i := 1; i <= card.quantity; i += 1 {
				c := Card{name: card.name, quantity: i}
				_, exists := metalist[c]
				if exists {
					metalist[c].count += 1
					metalist[c].total += card.quantity
				} else {
					metalist[c] = &Metacard{name: card.name, instance: i, total: card.quantity, count: 1}
				}
			}
		}
	}
	md := Metadeck{make([]*Metacard, 0)}
	for _, mc := range metalist {
		md.cards = append(md.cards, mc)
	}
	if len(md.cards) < decksize {
		return nil, errors.New(fmt.Sprintf("Cannot build %d card deck from %d cards", decksize, len(md.cards)))
	}
	sort.Sort(md)
	cards := make([]*Card, 0)
	for i := 0; i < decksize; i += 1 {
		mc := md.cards[i]
		found := false
		for _, card := range cards {
			if card.name == mc.name {
				card.quantity += 1
				found = true
				break
			}
		}
		if !found {
			cards = append(cards, &Card{name: mc.name, quantity: 1})
		}
	}
	return cards, nil
}

type Metadeck struct {
	cards []*Metacard
}

func (deck Metadeck) Len() int {
	return len(deck.cards)
}

func (deck Metadeck) Less(i, j int) bool {
	a := deck.cards[i]
	b := deck.cards[j]
	if a.count == b.count {
		if a.total == b.total {
			if a.name == b.name {
				return a.instance < b.instance
			} else {
				return a.name < b.name
			}
		} else {
			return a.total > b.total
		}
	} else {
		return a.count > b.count
	}
}

func (deck Metadeck) Swap(i, j int) {
	deck.cards[i], deck.cards[j] = deck.cards[j], deck.cards[i]
}
