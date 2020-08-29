package quotes_test

import (
	"github.com/ale-cci/gopm/quotes"
	"testing"
)

func TestFileIterator(t *testing.T) {
	t.Run("Get current quote", func(t *testing.T) {
		quote := quotes.PlainText{Text: "T", Author: "A"}
		chunks := []quotes.PlainText{quote}
		qi := quotes.NewFileIterator(chunks)

		current := qi.Current()

		if current != quote {
			t.Errorf("Wrong quote returned: %v expected %v", current, quote)
		}
	})

	t.Run("Should go to next quote", func(t *testing.T) {
		next := quotes.PlainText{Text: "Second quote"}
		chunks := []quotes.PlainText{{Text: "First quote"}, next}

		qi := quotes.NewFileIterator(chunks)

		qi.Next()
		current := qi.Current()

		if current != next {
			t.Errorf("Wrong quote returned: %v expected %v", current, next)
		}
	})

	t.Run("Should go to previous quote", func(t *testing.T) {
		first := quotes.PlainText{Text: "First quote"}
		chunks := []quotes.PlainText{first, {Text: "Second quote"}}

		qi := quotes.NewFileIterator(chunks)

		qi.Next()
		qi.Prev()
		current := qi.Current()

		if current != first {
			t.Errorf("Wrong quote returned: %v expected %v", current, first)
		}
	})

	t.Run("Should loop around on next", func(t *testing.T) {
		first := quotes.PlainText{Text: "First quote"}
		chunks := []quotes.PlainText{first, {Text: "Second quote"}, {Text: "Third quote"}}

		qi := quotes.NewFileIterator(chunks)

		qi.Next()
		qi.Next()
		qi.Next()
		current := qi.Current()

		if current != first {
			t.Errorf("Wrong quote returned: %v expected %v", current, first)
		}
	})

	t.Run("Should loop around on prev", func(t *testing.T) {
		first := quotes.PlainText{Text: "First quote"}
		chunks := []quotes.PlainText{first, {Text: "Second quote"}, {Text: "Third quote"}}

		qi := quotes.NewFileIterator(chunks)

		qi.Prev()
		qi.Prev()
		qi.Prev()
		current := qi.Current()

		if current != first {
			t.Errorf("Wrong quote returned: %v expected %v", current, first)
		}
	})

	t.Run("Should go to last quote on prev", func(t *testing.T) {
		last := quotes.PlainText{Text: "Last"}
		chunks := []quotes.PlainText{{Text: "First"}, last}

		fi := quotes.NewFileIterator(chunks)
		fi.Prev()

		current := fi.Current()

		if current != last {
			t.Errorf("Wrong quote returned: %v expected %v", current, last)
		}
	})
}
