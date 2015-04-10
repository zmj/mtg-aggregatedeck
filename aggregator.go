package main

import "regexp"
import "fmt"
import "errors"
import "strconv"
import "strings"
import "sort"

type Card struct {
	name string
	quantity int
}

type Deck struct {
	maindeck []*Card
	sideboard []*Card
}

var linePattern = regexp.MustCompile(`^(\d+) (.+)$`)
var maindecksize = 60
var sideboardsize = 15

func NewCard(raw string) (*Card,error) {
	segments := linePattern.FindStringSubmatch(raw)
	if len(segments) == 3 {
		count,_ := strconv.Atoi(segments[1])
		return &Card{ name:segments[2], quantity:count }, nil
	} else {
		return &Card{}, errors.New(fmt.Sprintf("Failed to parse line '%s'", raw))
	}
}

func NewDeck(raw []string) (*Deck,error) {
	maindeck := make([]*Card, 0)
	sideboard := make([]*Card, 0)
	count := 0
	for _,line := range raw {
		line = strings.TrimSpace(line)
		if len(line) == 0 || line=="Sideboard" {
			continue
		}
		card, err := NewCard(line)
		if err != nil {
			return nil, err
		}
		if count >= maindecksize {
			sideboard = append(sideboard, card)
		} else {
			maindeck = append(maindeck, card)
			count += card.quantity
		}
	}
	return &Deck{ maindeck: maindeck, sideboard: sideboard }, nil
}
	
func (deck *Deck) String() string {
	return "quack"
}

func aggregate(decks []*Deck) *Deck {
	maindecks := make([][]*Card, len(decks))
	sideboards := make([][]*Card, len(decks))
	for i,deck := range decks {
		maindecks[i] = deck.maindeck
		sideboards[i] = deck.sideboard
	}
	deck := &Deck { maindeck: agg(maindecks, maindecksize), sideboard: agg(sideboards, sideboardsize) }
	for _,c := range deck.maindeck {
		fmt.Printf("%s %d\n", c.name, c.quantity)
	}
	return deck
}

type Metacard struct {
	name string
	instance int	
	total int
	count int
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
	for i,mc := range md.cards {
		fmt.Printf("\t%d %d: %s %d %d\n", i, mc.count, mc.name, mc.instance, mc.total)
	}
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