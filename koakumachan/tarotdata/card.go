package tarotdata

import (
	"encoding/json"
	"os"
)

type Card struct {
	Name          string `json:"name"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	ImageFilename string `json:"image_filename"`

	// ImageData is the raw image data of the card.
	ImageData []byte
}

type Deck struct {
	MajorArcana []*Card `json:"major_arcana"`
	MinorArcana []*Card `json:"minor_arcana"`
	Court       []*Card `json:"court"`
}

func (d *Deck) Find(name string) *Card {
	for _, card := range d.MajorArcana {
		if card.Name == name {
			return card
		}
	}

	for _, card := range d.MinorArcana {
		if card.Name == name {
			return card
		}
	}

	for _, card := range d.Court {
		if card.Name == name {
			return card
		}
	}

	return nil
}

func LoadDeck(filename string) (*Deck, error) {
	bs, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var deck Deck
	err = json.Unmarshal(bs, &deck)
	if err != nil {
		return nil, err
	}

	for _, card := range deck.MajorArcana {
		imageData, err := os.ReadFile(card.ImageFilename)
		if err != nil {
			return nil, err
		}

		card.ImageData = imageData
	}

	for _, card := range deck.MinorArcana {
		imageData, err := os.ReadFile(card.ImageFilename)
		if err != nil {
			return nil, err
		}

		card.ImageData = imageData
	}

	for _, card := range deck.Court {
		imageData, err := os.ReadFile(card.ImageFilename)
		if err != nil {
			return nil, err
		}

		card.ImageData = imageData
	}

	return &deck, nil
}
