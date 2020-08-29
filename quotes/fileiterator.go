package quotes

type FileIterator struct {
	index  int
	chunks []PlainText
}

func NewFileIterator(chunks []PlainText) *FileIterator {
	return &FileIterator{index: 0, chunks: chunks}
}

func (qi *FileIterator) Current() PlainText {
	qi.normalizeIndex()
	return qi.chunks[qi.index]
}

func (qi *FileIterator) Next() {
	qi.index += 1
}

func (qi *FileIterator) Prev() {
	qi.index -= 1
}

func (qi *FileIterator) normalizeIndex() {
	index := qi.index % len(qi.chunks)
	if index < 0 {
		index += len(qi.chunks)
	}
	qi.index = index
}
