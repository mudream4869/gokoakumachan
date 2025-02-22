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
}

type Deck struct {
	MajorArcana []Card `json:"major_arcana"`
	MinorArcana []Card `json:"minor_arcana"`
	Court       []Card `json:"court"`
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

	return &deck, nil
}
