package prompt_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/vaultpull/internal/prompt"
)

func TestNew_DefaultsToStdStreams(t *testing.T) {
	c := prompt.New(nil, nil)
	if c == nil {
		t.Fatal("expected non-nil Confirmer")
	}
}

func TestConfirm_YesResponse(t *testing.T) {
	for _, input := range []string{"y", "Y", "yes", "YES"} {
		in := strings.NewReader(input + "\n")
		out := &bytes.Buffer{}
		c := prompt.New(in, out)

		ok, err := c.Confirm("Continue?", false)
		if err != nil {
			t.Fatalf("input %q: unexpected error: %v", input, err)
		}
		if !ok {
			t.Errorf("input %q: expected true, got false", input)
		}
	}
}

func TestConfirm_NoResponse(t *testing.T) {
	for _, input := range []string{"n", "N", "no", "NO"} {
		in := strings.NewReader(input + "\n")
		out := &bytes.Buffer{}
		c := prompt.New(in, out)

		ok, err := c.Confirm("Continue?", true)
		if err != nil {
			t.Fatalf("input %q: unexpected error: %v", input, err)
		}
		if ok {
			t.Errorf("input %q: expected false, got true", input)
		}
	}
}

func TestConfirm_EmptyInput_DefaultNo(t *testing.T) {
	in := strings.NewReader("\n")
	out := &bytes.Buffer{}
	c := prompt.New(in, out)

	ok, err := c.Confirm("Continue?", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected false for empty input with defaultYes=false")
	}
}

func TestConfirm_EmptyInput_DefaultYes(t *testing.T) {
	in := strings.NewReader("\n")
	out := &bytes.Buffer{}
	c := prompt.New(in, out)

	ok, err := c.Confirm("Continue?", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected true for empty input with defaultYes=true")
	}
}

func TestConfirm_InvalidResponse(t *testing.T) {
	in := strings.NewReader("maybe\n")
	out := &bytes.Buffer{}
	c := prompt.New(in, out)

	_, err := c.Confirm("Continue?", false)
	if err == nil {
		t.Fatal("expected error for invalid response")
	}
}

func TestConfirm_WritesPromptToOutput(t *testing.T) {
	in := strings.NewReader("y\n")
	out := &bytes.Buffer{}
	c := prompt.New(in, out)

	_, _ = c.Confirm("Apply changes?", false)

	if !strings.Contains(out.String(), "Apply changes?") {
		t.Errorf("expected prompt text in output, got: %q", out.String())
	}
	if !strings.Contains(out.String(), "y/N") {
		t.Errorf("expected hint [y/N] in output, got: %q", out.String())
	}
}
