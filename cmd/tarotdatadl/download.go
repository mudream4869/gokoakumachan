package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/jomei/notionapi"
	"github.com/mudream4869/gokoakumachan/koakumachan/kautil"
	"github.com/mudream4869/gokoakumachan/koakumachan/tarotdata"
)

const TAROT_IMAGES_DIR = "tarot_images"

const CARD_TYPE_MAJOR_ARCANA = "Major Arcana"
const CARD_TYPE_MINOR_ARCANA = "Minor Arcana"
const CARD_TYPE_COURT = "Court"

func getDataFromNotion(token, databaseID string) ([]notionapi.Page, error) {
	client := notionapi.NewClient(notionapi.Token(token))

	resp, err := client.Database.Query(
		context.Background(), notionapi.DatabaseID(databaseID), nil)

	if err != nil {
		return nil, kautil.Errorf("%w", err)
	}

	return resp.Results, nil
}

func getExtensionFromFilename(filename string) string {
	tmp := strings.Split(filename, ".")
	return tmp[len(tmp)-1]
}

func downloadFileFromURL(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return kautil.Errorf("%w", err)
	}
	defer resp.Body.Close()

	f, err := os.Create(filepath)
	if err != nil {
		return kautil.Errorf("%w", err)
	}

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return kautil.Errorf("%w", err)
	}

	return nil
}

type TarotData struct {
	MajorArcana []*tarotdata.Card `json:"major_arcana"`
	MinorArcana []*tarotdata.Card `json:"minor_arcana"`
	Court       []*tarotdata.Card `json:"court"`
}

func Main(confFilename string) error {
	conf, err := loadConf(confFilename)
	if err != nil {
		return kautil.Errorf("%w", err)
	}

	cred, err := loadCredential(conf.Notion.CredentialFile)
	if err != nil {
		return kautil.Errorf("%w", err)
	}

	data, err := getDataFromNotion(cred.Token, conf.Notion.TarotDataDatabaseID)
	if err != nil {
		return kautil.Errorf("%w", err)
	}

	os.Mkdir(TAROT_IMAGES_DIR, 0755)

	var fdata TarotData

	for _, row := range data {
		name := row.Properties["name"].(*notionapi.TitleProperty).Title[0].PlainText
		title := row.Properties["title"].(*notionapi.RichTextProperty).RichText[0].PlainText
		description := row.Properties["description"].(*notionapi.RichTextProperty).RichText[0].PlainText
		cardType := row.Properties["card_type"].(*notionapi.SelectProperty).Select.Name
		innerNumber := row.Properties["inner_number"].(*notionapi.NumberProperty).Number
		element := row.Properties["element"].(*notionapi.SelectProperty).Select.Name

		imgFiles := row.Properties["image"].(*notionapi.FilesProperty).Files

		imgFilename := ""
		if len(imgFiles) > 0 {
			imgFile := imgFiles[0]
			ext := getExtensionFromFilename(imgFile.Name)
			imgFilename = path.Join(TAROT_IMAGES_DIR, fmt.Sprintf("%s.%s", name, ext))

			url := imgFile.File.URL
			err := downloadFileFromURL(url, imgFilename)
			if err != nil {
				log.Println(name, err)
				imgFilename = ""
			}
		}

		card := tarotdata.Card{
			Name:          name,
			Title:         title,
			Description:   description,
			ImageFilename: imgFilename,
			Element:       element,
			InnerNumber:   int(innerNumber),
		}

		switch cardType {
		case CARD_TYPE_MAJOR_ARCANA:
			fdata.MajorArcana = append(fdata.MajorArcana, &card)
		case CARD_TYPE_MINOR_ARCANA:
			fdata.MinorArcana = append(fdata.MinorArcana, &card)
		case CARD_TYPE_COURT:
			fdata.Court = append(fdata.Court, &card)
		}
	}

	log.Println(len(fdata.MajorArcana), len(fdata.MinorArcana), len(fdata.Court))

	bs, err := json.MarshalIndent(fdata, "", "  ")
	if err != nil {
		return kautil.Errorf("%w", err)
	}

	err = os.WriteFile(conf.TarotDataFilename, bs, 0644)
	if err != nil {
		return kautil.Errorf("%w", err)
	}

	return nil
}

func main() {
	confFilename := flag.String("conf", "config.yaml", "config file name")
	flag.Parse()

	err := Main(*confFilename)
	if err != nil {
		log.Fatalln(err)
	}
}
