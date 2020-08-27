package quotes

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type Quote struct {
	Text   string `json:"text"`
	Author string `json:"author"`
}

func LoadFromJson(file io.Reader) ([]Quote, error) {
	bytes, _ := ioutil.ReadAll(file)

	var quotes []Quote
	err := json.Unmarshal([]byte(bytes), &quotes)

	if err != nil {
		return nil, err
	}

	return quotes, nil
}

func LoadFile(file io.Reader) (*Quote, error) {
	bytes, err := ioutil.ReadAll(file)

	if err != nil {
		return nil, err
	}

	return &Quote{Text: string(bytes)}, nil
}
