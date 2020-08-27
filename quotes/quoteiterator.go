package quotes

type QuoteIterator struct {
	index  int
	quotes []Quote
}

func NewQuoteIterator(quotes []Quote) *QuoteIterator {
	return &QuoteIterator{index: 0, quotes: quotes}
}

func (qi *QuoteIterator) Current() Quote {
	qi.normalizeIndex()
	return qi.quotes[qi.index]
}

func (qi *QuoteIterator) Next() {
	qi.index += 1
}

func (qi *QuoteIterator) Prev() {
	qi.index -= 1
}

func (qi *QuoteIterator) normalizeIndex() {
	index := qi.index % len(qi.quotes)
	if index < 0 {
		index += len(qi.quotes)
	}
	qi.index = index
}
