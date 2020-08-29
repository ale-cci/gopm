package chunk_test

import (
	"testing"

	"github.com/ale-cci/gopm/chunk"
	f "github.com/mattetti/filebuffer"
)

func TestChunkIterator(t *testing.T) {
	tt := []struct {
		name    string
		content []string
		next    int
		prev    int
		lines   int
		expect  string
	}{
		{
			name:    "Happy path",
			content: []string{"test"},
			next:    0,
			prev:    0,
			lines:   1,
			expect:  "test",
		},
		{
			name:    "next chunk in same file",
			content: []string{"first\nsecond"},
			next:    1,
			prev:    0,
			lines:   1,
			expect:  "second",
		},
		{
			name:    "next chunk on second file",
			content: []string{"first", "second"},
			next:    1,
			prev:    0,
			lines:   1,
			expect:  "second",
		},
		{
			name:    "Should loop over files",
			content: []string{"first", "second"},
			next:    2,
			prev:    0,
			lines:   1,
			expect:  "first",
		},
		{
			name:    "Should return on previous file on prev",
			content: []string{"first", "second", "third"},
			next:    1,
			prev:    1,
			lines:   1,
			expect:  "first",
		},
		{
			name:    "Should return to previous quote",
			content: []string{"first\nfirst2", "second", "third"},
			next:    1,
			prev:    1,
			lines:   1,
			expect:  "first",
		},
		{
			name:    "Should not skip last quote of each file",
			content: []string{"first\nfirst2", "second", "third"},
			next:    1,
			prev:    0,
			lines:   1,
			expect:  "first2",
		},
		{
			name:    "Should start file from zero if it's unvisited",
			content: []string{"first\nfirst2", "second", "third\nthird2"},
			next:    0,
			prev:    1,
			lines:   1,
			expect:  "third",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			chunkedFiles := make([]chunk.Iterator, len(tc.content))

			for i, text := range tc.content {
				file := f.New([]byte(text))
				chunkedFiles[i] = chunk.NewChunkedFile(file, tc.lines)
			}

			ci := chunk.ChunkIterator{Files: chunkedFiles}

			for i := 0; i < tc.next; i++ {
				ci.Next()
			}
			for i := 0; i < tc.prev; i++ {
				ci.Prev()
			}

			got := ci.Current()
			expect := tc.expect

			if got != expect {
				t.Fatalf("Wrong chunk returned: %q, expected %q", got, expect)
			}
		})
	}
}
