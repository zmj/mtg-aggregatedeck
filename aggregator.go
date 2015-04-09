package main

import "regexp"
import "fmt"
import "errors"
import "strconv"
import "strings"

type Card struct {
	name string
	quantity int
}

type Deck struct {
	maindeck []Card
	sideboard []Card
}

var linePattern = regexp.MustCompile(`^(\d+) (.+)$`)

func NewCard(raw string) (Card,error) {
	segments := linePattern.FindStringSubmatch(raw)
	if len(segments) == 3 {
		count,_ := strconv.Atoi(segments[1])
		return Card{ name:segments[2], quantity:count }, nil
	} else {
		return Card{}, errors.New(fmt.Sprintf("Failed to parse line '%s'", raw))
	}
}

func NewDeck(raw []string) (*Deck,error) {
	maindeck := make([]Card, 0)
	sideboard := make([]Card, 0)
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
		if count >= 60 {
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
	for _,card := range decks[0].sideboard {
		fmt.Printf("%d %s\n", card.quantity, card.name)
	}
	return decks[0]
}