package quotes

import (
	"strings"
	"testing"
)

func TestLoadFromJson(t *testing.T) {
	t.Run("Happy path", func(t *testing.T) {
		data := strings.NewReader(`[
			{"text": "First quote text", "author": "test-author"},
			{"text": "Second quote text", "author": "second-author"}
		]`)

		quotes, err := LoadFromJson(data)
		if err != nil {
			t.Error(err)
		}

		if len(quotes) != 2 {
			t.Errorf("Quote count %d != 2", len(quotes))
		}

		if quotes[0].Author != "test-author" {
			t.Errorf("Mismatch author: '%s' expected 'test-author'", quotes[0].Author)
		}

		if quotes[0].Text != "First quote text" {
			t.Errorf("Mismatch text: '%s', expected 'First quote text'", quotes[0].Text)
		}
	})

	t.Run("Returns error when json file is malformed", func(t *testing.T) {
		malformed := strings.NewReader("{}")
		quotes, err := LoadFromJson(malformed)

		if quotes != nil {
			t.Errorf("List of quotes != nil: %q", quotes)
		}

		if err == nil {
			t.Errorf("Error == nil")
		}
	})
}

func TestQuoteIterator(t *testing.T) {
	t.Run("Get current quote", func(t *testing.T) {
		quote := Quote{Text: "T", Author: "A"}
		quotes := []Quote{quote}
		qi := NewQuoteIterator(quotes)

		current := qi.Current()

		if current != quote {
			t.Errorf("Wrong quote returned: %v expected %v", current, quote)
		}
	})

	t.Run("Should go to next quote", func(t *testing.T) {
		next := Quote{Text: "Second quote"}
		quotes := []Quote{{Text: "First quote"}, next}

		qi := NewQuoteIterator(quotes)

		qi.Next()
		current := qi.Current()

		if current != next {
			t.Errorf("Wrong quote returned: %v expected %v", current, next)
		}
	})

	t.Run("Should go to previous quote", func(t *testing.T) {
		first := Quote{Text: "First quote"}
		quotes := []Quote{first, {Text: "Second quote"}}

		qi := NewQuoteIterator(quotes)

		qi.Next()
		qi.Prev()
		current := qi.Current()

		if current != first {
			t.Errorf("Wrong quote returned: %v expected %v", current, first)
		}
	})

	t.Run("Should loop around on next", func(t *testing.T) {
		first := Quote{Text: "First quote"}
		quotes := []Quote{first, {Text: "Second quote"}, {Text: "Third quote"}}

		qi := NewQuoteIterator(quotes)

		qi.Next()
		qi.Next()
		qi.Next()
		current := qi.Current()

		if current != first {
			t.Errorf("Wrong quote returned: %v expected %v", current, first)
		}
	})

	t.Run("Should loop around on prev", func(t *testing.T) {
		first := Quote{Text: "First quote"}
		quotes := []Quote{first, {Text: "Second quote"}, {Text: "Third quote"}}

		qi := NewQuoteIterator(quotes)

		qi.Prev()
		qi.Prev()
		qi.Prev()
		current := qi.Current()

		if current != first {
			t.Errorf("Wrong quote returned: %v expected %v", current, first)
		}
	})
}
