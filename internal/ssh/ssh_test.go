package ssh

import (
	"testing"

	"github.com/byawitz/ggh/internal/config"
)

func TestGenerateCommandArgs_SetEnv(t *testing.T) {
	c := config.SSHConfig{
		Host:   "example.com",
		User:   "admin",
		Port:   "22",
		Key:    "~/.ssh/id_rsa",
		SetEnv: []string{"TERM=xterm-256color", "COLORTERM=truecolor"},
	}

	args := GenerateCommandArgs(c)

	// Check basic args
	if args[0] != "admin@example.com" {
		t.Errorf("Expected admin@example.com, got %s", args[0])
	}

	// Collect all -o args
	var opts []string
	for i, a := range args {
		if a == "-o" && i+1 < len(args) {
			opts = append(opts, args[i+1])
		}
	}

	// Should have both SetEnv values, quoted
	expected := []string{
		"'SetEnv=TERM=xterm-256color'",
		"'SetEnv=COLORTERM=truecolor'",
	}
	for _, exp := range expected {
		found := false
		for _, opt := range opts {
			if opt == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected option %s not found in args: %v", exp, args)
		}
	}
}

func TestGenerateCommandArgs_NoSetEnv(t *testing.T) {
	c := config.SSHConfig{
		Host: "example.com",
		User: "root",
	}

	args := GenerateCommandArgs(c)

	for _, a := range args {
		if a == "'SetEnv='" {
			t.Error("Empty SetEnv should not appear in args")
		}
	}
}
