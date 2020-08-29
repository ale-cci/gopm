package chunk

import (
	"bufio"
	"io"
	"strings"
)

type CircularItarator interface {
	Next()
	Prev()
	Current() []rune
}

type ChunkIterator struct {
	Files []io.ReadSeeker
	Lines int

	// Index representing current file
	index int
}

func (c *ChunkIterator) Next() {

}

func (c *ChunkIterator) Prev() {
}

func (c *ChunkIterator) Current() ([]rune, error) {
	file := c.Files[c.index]

	buffer := ""
	reader := bufio.NewReader(file)

	for i := 0; i < c.Lines; i++ {
		line, err := reader.ReadString('\n')

		if err != nil && err != io.EOF {
			return nil, err
		}

		buffer += string(line)
	}

	// Remove trailing zeroes
	buffer = strings.TrimRight(buffer, "\n")
	return []rune(buffer), nil
}
