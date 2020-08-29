package chunk

import (
	"bufio"
	"io"
	"strings"
)

type ChunkedFile struct {
	File   io.ReadSeeker
	lines  int
	index  int
	chunks []int64
	eof    bool
}

func NewChunkedFile(file io.ReadSeeker, lines int) *ChunkedFile {
	return &ChunkedFile{File: file, lines: lines, index: 0, chunks: []int64{0}}
}

func (cf *ChunkedFile) isLastChunk() bool {
	return cf.index+1 == len(cf.chunks)
}

func (cf *ChunkedFile) Next() bool {
	// Calculate chunk size
	if cf.isLastChunk() && !cf.eof {
		currentChunk := cf.chunks[cf.index]
		buffer, eof := cf.readLines()

		cf.eof = eof

		if !cf.eof {
			currentChunk += int64(len(buffer))
			cf.chunks = append(cf.chunks, currentChunk)
		}
	}

	if !cf.isLastChunk() {
		cf.index++
		return cf.eof
	} else {
		return true
	}
}

func (cf *ChunkedFile) Prev() bool {
	if cf.index == 0 {
		return true
	}

	cf.index--
	return false
}

func (cf *ChunkedFile) readLines() (string, bool) {
	cf.File.Seek(cf.chunks[cf.index], io.SeekStart)

	buffer := ""
	reader := bufio.NewReader(cf.File)

	for i := 0; i < cf.lines; i++ {
		line, err := reader.ReadString('\n')
		buffer += line

		if err != nil {
			if err == io.EOF {
				return buffer, true
			} else {
				panic(err)
			}
		}
	}

	return buffer, false
}

func (cf *ChunkedFile) Current() string {
	buffer, _ := cf.readLines()
	buffer = strings.TrimRight(buffer, "\n")

	return buffer
}
