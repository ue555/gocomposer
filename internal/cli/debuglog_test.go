package cli

import (
	"bytes"
	"strings"
	"testing"

	"gocomposer/internal/cio"
	"gocomposer/internal/config"
)

func TestDeprecationNotices(t *testing.T) {
	defer func() { config.Deprecation = nil }()
	var buf bytes.Buffer
	io := &cio.IO{Out: &bytes.Buffer{}, Err: &buf, Verbosity: cio.Normal}
	registerErrorHandler(io)
	for i := 0; i < 3; i++ {
		config.Deprecation("something is deprecated.")
	}
	want := "Deprecation Notice: something is deprecated. in gocomposer\nMore deprecation notices were hidden, run again with `-v` to show them.\n"
	if buf.String() != want {
		t.Errorf("output = %q", buf.String())
	}
	buf.Reset()
	io.Verbosity = cio.Verbose
	registerErrorHandler(io)
	config.Deprecation("a")
	config.Deprecation("b")
	if strings.Count(buf.String(), "Deprecation Notice:") != 2 {
		t.Errorf("verbose output = %q", buf.String())
	}
}

func TestBinCompatDeprecation(t *testing.T) {
	defer func() { config.Deprecation = nil }()
	var got []string
	config.Deprecation = func(m string) { got = append(got, m) }
	c := config.New(false, "/tmp")
	c.Set("bin-compat", "symlink")
	if _, err := c.GetFlags("bin-compat", 0); err != nil {
		t.Fatal(err)
	}
	c.Set("bin-compat", "proxy")
	_, _ = c.GetFlags("bin-compat", 0)
	if len(got) != 1 || !strings.Contains(got[0], `config.bin-compat "symlink" is deprecated`) {
		t.Errorf("notices = %q", got)
	}
}

func TestDebugVerbosity(t *testing.T) {
	for _, c := range []struct {
		argv []string
		want int
	}{
		{[]string{"show", "-vvv"}, cio.Debug},
		{[]string{"show", "-q", "-vvv"}, cio.Quiet},
		{[]string{"show", "--verbose=3"}, cio.Debug},
		{[]string{"show"}, cio.Normal},
	} {
		if got := debugVerbosity(c.argv); got != c.want {
			t.Errorf("%v: %d", c.argv, got)
		}
	}
}
