package main

import (
	"os"

	"github.com/mudream4869/gokoakumachan/koakumachan/kautil"
	"gopkg.in/yaml.v3"
)

type NotionCredential struct {
	Token string `yaml:"token"`
}

type Notion struct {
	CredentialFile      string `yaml:"credential_file"`
	TarotDataDatabaseID string `yaml:"tarot_data_database_id"`
}

type Conf struct {
	TarotDataFilename string `yaml:"tarot_data_filename"`
	Notion            Notion `yaml:"notion"`
}

func loadConf(confFilename string) (*Conf, error) {
	bs, err := os.ReadFile(confFilename)
	if err != nil {
		return nil, kautil.Errorf("%w", err)
	}

	var conf Conf
	err = yaml.Unmarshal(bs, &conf)
	if err != nil {
		return nil, kautil.Errorf("%w", err)
	}

	return &conf, nil
}

func loadCredential(credentialFile string) (*NotionCredential, error) {
	bs, err := os.ReadFile(credentialFile)
	if err != nil {
		return nil, kautil.Errorf("%w", err)
	}

	var cred NotionCredential
	err = yaml.Unmarshal(bs, &cred)
	if err != nil {
		return nil, kautil.Errorf("%w", err)
	}

	return &cred, nil
}
