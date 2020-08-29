package quotes

import (
	"io"
	"io/ioutil"
)

type PlainText struct {
	Text   string
	Author string
}

func LoadFile(file io.Reader) (*PlainText, error) {
	bytes, err := ioutil.ReadAll(file)

	if err != nil {
		return nil, err
	}

	return &PlainText{Text: string(bytes)}, nil
}
