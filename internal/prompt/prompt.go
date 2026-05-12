package prompt

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Confirmer handles interactive user confirmation prompts.
type Confirmer struct {
	in  io.Reader
	out io.Writer
}

// New creates a new Confirmer with the given input and output streams.
// If nil is provided, os.Stdin and os.Stdout are used.
func New(in io.Reader, out io.Writer) *Confirmer {
	if in == nil {
		in = os.Stdin
	}
	if out == nil {
		out = os.Stdout
	}
	return &Confirmer{in: in, out: out}
}

// Confirm displays a yes/no prompt with the given message and returns
// true if the user confirms, false otherwise. Defaults to false on
// empty input unless defaultYes is true.
func (c *Confirmer) Confirm(message string, defaultYes bool) (bool, error) {
	defaultHint := "y/N"
	if defaultYes {
		defaultHint = "Y/n"
	}

	_, err := fmt.Fprintf(c.out, "%s [%s]: ", message, defaultHint)
	if err != nil {
		return false, fmt.Errorf("prompt write: %w", err)
	}

	var response string
	buf := make([]byte, 256)
	n, err := c.in.Read(buf)
	if err != nil && err != io.EOF {
		return false, fmt.Errorf("prompt read: %w", err)
	}
	response = strings.TrimSpace(string(buf[:n]))

	if response == "" {
		return defaultYes, nil
	}

	switch strings.ToLower(response) {
	case "y", "yes":
		return true, nil
	case "n", "no":
		return false, nil
	default:
		return false, fmt.Errorf("unrecognised response %q: expected y/yes or n/no", response)
	}
}
