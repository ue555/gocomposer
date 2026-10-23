package cio

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func qio(input string) (*IO, *bytes.Buffer) {
	var errb bytes.Buffer
	return &IO{Out: &bytes.Buffer{}, Err: &errb, Verbosity: Normal, Interactive: true, In: strings.NewReader(input)}, &errb
}

func TestAskConfirmationRetries(t *testing.T) {
	o, errb := qio("maybe\nN\n")
	ok, err := o.AskConfirmation("Sure? ", true)
	if err != nil || ok {
		t.Fatalf("got %v %v", ok, err)
	}
	want := "Sure? Please answer yes, y, no, or n.\nSure? "
	if errb.String() != want {
		t.Errorf("output %q", errb.String())
	}
	o, _ = qio("\n")
	if ok, _ := o.AskConfirmation("Sure? ", true); !ok {
		t.Error("empty answer must give the default")
	}
	o, _ = qio("0\n")
	if ok, _ := o.AskConfirmation("Sure? ", true); !ok {
		t.Error(`"0" is empty() and gives the true default`)
	}
}

func TestAskEOFFallsBackToDefault(t *testing.T) {
	o, _ := qio("")
	v, err := o.Ask("Name: ", "def")
	if err != nil || v != "def" || o.Interactive {
		t.Fatalf("got %v %v interactive=%v", v, err, o.Interactive)
	}
	o, _ = qio("")
	if _, err := o.Ask("Name: ", nil); !errors.As(err, new(*MissingInputError)) {
		t.Fatalf("want Aborted., got %v", err)
	}
}

func TestAskAndValidateAttempts(t *testing.T) {
	o, errb := qio("a\nb\nc\n")
	calls := 0
	_, err := o.AskAndValidate("Q: ", func(v any) (any, error) {
		calls++
		return nil, &InvalidArgumentError{Msg: "bad " + v.(string)}
	}, 2, nil)
	if err == nil || err.Error() != "bad b" || calls != 2 {
		t.Fatalf("got %v after %d calls", err, calls)
	}
	if errb.String() != "Q: bad a\nQ: " {
		t.Errorf("output %q", errb.String())
	}
}

func TestSelect(t *testing.T) {
	choices := []Choice{{"0", "alpha"}, {"1", "beta"}}
	o, errb := qio("1\n")
	v, err := o.Select("Pick", choices, false, nil, 0, `Value "%s" is invalid`, false)
	if err != nil || v != "1" {
		t.Fatalf("got %v %v", v, err)
	}
	if errb.String() != "Pick\n  [0] alpha\n  [1] beta\n > " {
		t.Errorf("output %q", errb.String())
	}
	o, _ = qio("beta,alpha\n")
	v, err = o.Select("Pick", choices, false, nil, 0, `Value "%s" is invalid`, true)
	if err != nil || strings.Join(v.([]string), ",") != "0,1" {
		t.Fatalf("multiselect got %v %v", v, err)
	}
	o, _ = qio("gamma\n")
	if _, err := o.Select("Pick", choices, false, nil, 1, `Value "%s" is invalid`, false); err == nil || err.Error() != `Value "gamma" is invalid` {
		t.Fatalf("got %v", err)
	}
	assoc := []Choice{{"build", "Build it"}, {"test", ""}}
	o, _ = qio("test\n")
	if v, _ := o.Select("Script", assoc, true, "", 1, `Invalid script name "%s"`, false); v != "test" {
		t.Fatalf("assoc got %v", v)
	}
}

func TestNonInteractiveDefaults(t *testing.T) {
	o := &IO{Out: &bytes.Buffer{}, Err: &bytes.Buffer{}, Verbosity: Normal}
	v, err := o.AskAndValidate("Q: ", func(v any) (any, error) { return v.(string) + "!", nil }, 0, "x")
	if err != nil || v != "x!" {
		t.Fatalf("got %v %v", v, err)
	}
	if v, _ := o.Ask("Q: ", nil); v != nil {
		t.Fatalf("got %v", v)
	}
}

func TestStatusEncoding(t *testing.T) {
	st := &ProgramStatus{State: "blocked", Message: "Do you <info>want</info>?\x07", Kind: KindPermission, hasMessage: true}
	got := encodeStatus("7501", st)
	if got != "state=blocked:app=composer:kind=permission:msg=RG8geW91IHdhbnQ/" {
		t.Errorf("got %q", got)
	}
	if encodeStatus("9;4", &ProgramStatus{State: "working"}) != "3" {
		t.Error("9;4 working")
	}
}
