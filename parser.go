package main

import "regexp"
import "fmt"
import "errors"
import "strconv"
import "strings"
import "bytes"

type Card struct {
	name string
	quantity int
}

type Deck struct {
	maindeck []*Card
	sideboard []*Card
}

var linePattern = regexp.MustCompile(`^[^0-9]*(?P<quantity>\d+) (\[.+\] )?(?P<name>.+)$`)

func NewCard(line string) (*Card,error) {
	match := linePattern.FindStringSubmatch(line)
	card := &Card{ }
	if match==nil {
		return card, errors.New(fmt.Sprintf("Failed to parse line '%v'", []byte(line)))
	}
	for i,group := range linePattern.SubexpNames() {	
		if group=="quantity" {
			quantity, err := strconv.Atoi(match[i])
			if err != nil {
				return card, err
			}
			card.quantity = quantity
		} else if group=="name" {
			card.name = match[i]
		}
	}
	if len(card.name)==0 {
		return nil, errors.New(fmt.Sprintf("Could not parse card name from '%s'", line))
	} else if card.quantity==0 {
		return nil, errors.New(fmt.Sprintf("Could not parse card quantity from '%s'", line))
	} else {
		return card, nil
	}
}

func ignoreDecklistLine(line string) bool {
	if len(line) == 0 ||
		line == "Sideboard" ||
		strings.HasPrefix(line, "//") ||
		strings.HasPrefix(line, "#") {
		return true
	} 
	return false	
}

func NewDeck(raw []string) (*Deck,error) {
	maindeck := make([]*Card, 0)
	sideboard := make([]*Card, 0)
	count := 0
	for _,line := range raw {
		line = strings.TrimSpace(line)
		if ignoreDecklistLine(line) {
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
	s := bytes.Buffer{}
	for _,c := range deck.maindeck {
		s.WriteString(fmt.Sprintf("%d %s\n", c.quantity, c.name))
	}
	s.WriteString("\nSideboard\n")
	for _,c := range deck.sideboard {
		s.WriteString(fmt.Sprintf("%d %s\n", c.quantity, c.name))
	}
	return s.String()
}