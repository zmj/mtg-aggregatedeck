package main

type Deck struct {

}

func NewDeck(raw []string) Deck {
	return Deck {}
}
	
func (deck Deck) String() string {
	return "quack"
}

func aggregate(decks []Deck) Deck {
	return decks[0]
}