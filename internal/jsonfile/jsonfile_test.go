package jsonfile

import "testing"

func TestParseErrors(t *testing.T) {
	cases := []struct{ in, file, want string }{
		{"{\"a\": 1,}", "composer.json", "\"composer.json\" does not contain valid JSON\nParse error on line 1:\n{\"a\": 1,}\n-------^\nExpected: 'STRING' - It appears you have an extra trailing comma"},
		{"{\"a\": \"\xff\"}", "x.json", "\"x.json\" is not UTF-8, could not parse as JSON"},
		{"{\"a\": 1}", "", ""},
	}
	for _, c := range cases {
		_, err := Parse([]byte(c.in), c.file)
		got := ""
		if err != nil {
			got = err.Error()
		}
		if got != c.want {
			t.Errorf("%q: got %q want %q", c.in, got, c.want)
		}
	}
}

func TestLockConflictRecovery(t *testing.T) {
	in := "{\n<<<<<<< HEAD\n    \"content-hash\": \"abc123\",\n=======\n    \"content-hash\": \"def456\",\n>>>>>>> branch\n    \"packages\": []\n}\n"
	v, err := Parse([]byte(in), "composer.lock")
	if err != nil || v == nil {
		t.Fatalf("expected recovery, got %v", err)
	}
}
