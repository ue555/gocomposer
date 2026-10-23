package cio

import (
	"strings"
	"testing"
)

// Expectations verified against symfony/console Table rendered by the oracle
// PHP 8.4 image with Composer f83588392df43d4fd027d84d26de459946ef4e80 deps.
func TestTableDefaultStyle(t *testing.T) {
	tb := NewTable("default")
	tb.Headers = []string{"ISBN", "Title", "Author"}
	tb.Rows = [][]string{
		{"99921-58-10-7", "Divine Comedy", "Dante Alighieri"},
		{"9971-5-0210-0", "A Tale of Two Cities", "Charles Dickens"},
		{"960-425-059-0", "The Lord of the Rings", "J. R. R. Tolkien"},
		{"80-902734-1-6", "And Then There Were None", "Agatha Christie"},
	}
	var out []string
	tb.RenderTo(func(s string) { out = append(out, StripTags(s)) })
	want := `+---------------+--------------------------+------------------+
| ISBN          | Title                    | Author           |
+---------------+--------------------------+------------------+
| 99921-58-10-7 | Divine Comedy            | Dante Alighieri  |
| 9971-5-0210-0 | A Tale of Two Cities     | Charles Dickens  |
| 960-425-059-0 | The Lord of the Rings    | J. R. R. Tolkien |
| 80-902734-1-6 | And Then There Were None | Agatha Christie  |
+---------------+--------------------------+------------------+`
	if got := strings.Join(out, "\n"); got != want {
		t.Fatalf("got\n%s\nwant\n%s", got, want)
	}
}

func TestTableCompactStyle(t *testing.T) {
	tb := NewTable("compact")
	tb.Headers = []string{"ISBN", "Title", "Author"}
	tb.Rows = [][]string{
		{"99921-58-10-7", "Divine Comedy", "Dante Alighieri"},
		{"9971-5-0210-0", "A Tale of Two Cities", "Charles Dickens"},
		{"960-425-059-0", "The Lord of the Rings", "J. R. R. Tolkien"},
		{"80-902734-1-6", "And Then There Were None", "Agatha Christie"},
	}
	var out []string
	tb.RenderTo(func(s string) { out = append(out, StripTags(s)) })
	want := `ISBN          Title                    Author           
99921-58-10-7 Divine Comedy            Dante Alighieri  
9971-5-0210-0 A Tale of Two Cities     Charles Dickens  
960-425-059-0 The Lord of the Rings    J. R. R. Tolkien 
80-902734-1-6 And Then There Were None Agatha Christie  `
	if got := strings.Join(out, "\n"); got != want {
		t.Fatalf("got\n%q\nwant\n%q", got, want)
	}
}

func TestTableMultilineCell(t *testing.T) {
	tb := NewTable("default")
	tb.Rows = [][]string{{"a\nb", "c"}}
	var out []string
	tb.RenderTo(func(s string) { out = append(out, s) })
	want := "+---+---+\n| a | c |\n| b |   |\n+---+---+"
	if got := strings.Join(out, "\n"); got != want {
		t.Fatalf("got\n%s\nwant\n%s", got, want)
	}
}
