package main

import "sort"

var maindecksize = 60
var sideboardsize = 15

func aggregate(decks []*Deck) *Deck {
	maindecks := make([][]*Card, len(decks))
	sideboards := make([][]*Card, len(decks))
	for i,deck := range decks {
		maindecks[i] = deck.maindeck
		sideboards[i] = deck.sideboard
	}
	deck := &Deck { maindeck: agg(maindecks, maindecksize), sideboard: agg(sideboards, sideboardsize) }
	return deck
}

type Metacard struct {
	name string
	instance int // nth copy of card
	total int // total # of card played in all decklists
	count int // # of nth copy in all decklists
}

func agg(cardlists [][]*Card, decksize int) []*Card {
	metalist := map[Card]*Metacard{}
	for _,list := range cardlists {
		for _,card := range list {
			for i:=1; i<=card.quantity; i+=1 {
				c := Card{ name: card.name, quantity: i }
				_,exists := metalist[c]
				if exists {
					metalist[c].count += 1
					metalist[c].total += card.quantity
				} else {
					metalist[c] = &Metacard{ name: card.name, instance: i, total: card.quantity, count: 1 }
				}
			}
		}
	}
	md := Metadeck{ make([]*Metacard, 0) }
	for _,mc := range metalist {
		md.cards = append(md.cards, mc)
	}
	sort.Sort(md)
	cards := make([]*Card, 0)
	for i:=0; i<decksize; i+=1 {
		card := md.cards[i]
		found := false
		for _,c := range cards {
			if c.name == card.name {
				c.quantity += 1
				found = true
				break
			}
		}
		if !found {
			cards = append(cards, &Card{ name: card.name, quantity: 1 })
		}
	}
	return cards
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