package chunk

import (
	"io"
	"testing"

	f "github.com/mattetti/filebuffer"
)

func TestChunkIterator(t *testing.T) {
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
				ci := ChunkIterator{Files: []io.ReadSeeker{fb}, Lines: tc.linesToRead}

				runes, err := ci.Current()

				got := string(runes)
				expect := tc.expect

				if err != nil {
					t.Errorf("Error encountered in %q: %v", tc.name, err)
				}

				if got != expect {
					t.Fatalf("Wrong chunk: %q, expected: %q", got, expect)
				}
			})
		}
	})

}
