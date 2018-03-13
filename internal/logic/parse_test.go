package logic

import (
	"errors"
	"os"
	"testing"
)

func TestScgDecklist(t *testing.T) {
	filename := "scg_grag_1.txt"
	_, err := ParseDecklist(filename)
	if err != nil {
		t.Error(err)
	}
}

func TestWotcDecklist(t *testing.T) {
	filename := "wotc_ptdtk_abzag_1.txt"
	_, err := ParseDecklist(filename)
	if err != nil {
		t.Error(err)
	}
}

func TestMwsDecklist(t *testing.T) {
	filename := "mws_modabz_1.mwDeck"
	_, err := ParseDecklist(filename)
	if err != nil {
		t.Error(err)
	}
}

func ParseDecklist(filename string) (*Deck, error) {
	file, err := os.Open("sample/" + filename)
	if err != nil {
		return nil, err
	}
	deck, err := NewDeck(file)
	if err != nil {
		return nil, err
	}
	if CountCards(deck.maindeck) < 60 {
		return nil, errors.New("Maindeck less than 60 cards")
	}
	if CountCards(deck.sideboard) > 15 {
		return nil, errors.New("Sideboard more than 15 cards")
	}
	return deck, nil
}

func CountCards(deck []*Card) (count int) {
	for _, c := range deck {
		count += c.quantity
	}
	return count
}
