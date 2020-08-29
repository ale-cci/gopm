package chunk_test

import (
	"testing"

	"github.com/ale-cci/gopm/chunk"
	f "github.com/mattetti/filebuffer"
)

func TestChunkedFile(t *testing.T) {
	t.Run("Current()", func(t *testing.T) {
		tt := []struct {
			content     string
			linesToRead int
			expect      string
			name        string
		}{
			{
				content:     "test\n",
				linesToRead: 1,
				expect:      "test",
				name:        "Two lines read one",
			},
			{
				content:     "multiline\ntext",
				linesToRead: 2,
				expect:      "multiline\ntext",
				name:        "Two lines read two",
			},
			{
				content:     "short",
				linesToRead: 2,
				expect:      "short",
				name:        "One line read two",
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				fb := f.New([]byte(tc.content))
				ci := chunk.NewChunkedFile(fb, tc.linesToRead)

				got := ci.Current()
				expect := tc.expect

				if got != expect {
					t.Fatalf("Wrong chunk: %q, expected: %q", got, expect)
				}
			})
		}
	})

	t.Run("Current called twice should return the same list of runes", func(t *testing.T) {
		fb := f.New([]byte("line1\nline2"))
		ci := chunk.NewChunkedFile(fb, 1)

		ci.Current()
		got := ci.Current()
		expect := "line1"

		if got != expect {
			t.Errorf("Wrong chunk %q, expected: %q", got, expect)
		}
	})

	t.Run("Next()", func(t *testing.T) {
		tt := []struct {
			name     string
			content  string
			lines    int
			expect   string
			iter     int
			stopiter bool
		}{
			{
				name:     "Happy path",
				content:  "test\nmultiline\nstring",
				lines:    1,
				expect:   "multiline",
				iter:     1,
				stopiter: false,
			},
			{
				name:     "Next on two lines chunk",
				content:  "ab\nbaaa\nc\nd",
				lines:    2,
				expect:   "c\nd",
				iter:     1,
				stopiter: false,
			},
			{
				name:     "Chunks of different size",
				content:  "first\nsecond\nthird",
				lines:    1,
				expect:   "third",
				iter:     2,
				stopiter: false,
			},
			{
				name:     "Next until end of file",
				content:  "a\nb",
				lines:    1,
				expect:   "b",
				iter:     3,
				stopiter: true,
			},
			{
				name:     "Should interrupt when end is reached",
				content:  "a\nb",
				lines:    1,
				expect:   "b",
				iter:     2,
				stopiter: true,
			},
			{
				name:     "Should interrupt when end is reached",
				content:  "ab",
				lines:    1,
				expect:   "ab",
				iter:     2,
				stopiter: true,
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				fb := f.New([]byte(tc.content))
				ci := chunk.NewChunkedFile(fb, tc.lines)

				var stopiteration bool
				for i := 0; i < tc.iter; i++ {
					stopiteration = ci.Next()
				}

				if tc.stopiter != stopiteration {
					t.Errorf("On %q StopIteration has value of: %v, expected: %v", tc.name, stopiteration, tc.stopiter)
				}
				got := ci.Current()
				expect := tc.expect

				if got != expect {
					t.Fatalf("At %q, wrong chunk: %q, expected: %q", tc.name, got, expect)
				}
			})
		}
	})

	t.Run("Prev()", func(t *testing.T) {
		tt := []struct {
			name    string
			next    int
			prev    int
			content string
			si      bool
			expect  string
		}{
			{
				name:    "Happy path",
				content: "test\ntest",
				next:    1,
				prev:    1,
				si:      false,
				expect:  "test",
			},
			{
				name:    "stopiteration if prev invoked more than next",
				content: "first\nsecond\nthird",
				next:    1,
				prev:    2,
				si:      true,
				expect:  "first",
			},
			{
				name:    "return stopiteration if invoked at beginning",
				content: "test",
				next:    0,
				prev:    1,
				si:      true,
				expect:  "test",
			},
			{
				name:    "shuld return back to previous chunk",
				content: "first\nsecond\nthird",
				next:    2,
				prev:    1,
				si:      false,
				expect:  "second",
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				fb := f.New([]byte(tc.content))
				cf := chunk.NewChunkedFile(fb, 1)

				for i := 0; i < tc.next; i++ {
					cf.Next()
				}

				si := false
				for i := 0; i < tc.prev && !si; i++ {
					si = cf.Prev()
				}

				if si != tc.si {
					t.Fatalf("Unexpected StopIteration on %q value: %v, expected: %v", tc.name, si, tc.si)
				}

				chunk := cf.Current()

				got := string(chunk)
				expect := tc.expect

				if got != expect {
					t.Fatalf("Unexpected chunk for %q, got: %q, expect: %q", tc.name, got, expect)
				}
			})
		}
	})

	t.Run("Combination of prev and next", func(t *testing.T) {
		fb := f.New([]byte("first\nsecond\nthird"))
		cf := chunk.NewChunkedFile(fb, 1)

		cf.Next()
		cf.Prev()
		cf.Next()
		cf.Next()

		chunk := cf.Current()

		got := string(chunk)
		expect := "third"

		if got != expect {
			t.Fatalf("Unexpected chunk: %q, expected: %q", got, expect)
		}
	})
}
