package main

import (
	"strings"
	"testing"
)

func TestGenerateTokenLength(t *testing.T) {
	for _, length := range []int{16, 32, 64} {
		tok := generateToken(length)
		if len(tok) != length {
			t.Errorf("generateToken(%d): got length %d", length, len(tok))
		}
	}
}

func TestGenerateTokenUniqueness(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 20; i++ {
		tok := generateToken(32)
		if seen[tok] {
			t.Errorf("duplicate token generated: %q", tok)
		}
		seen[tok] = true
	}
}

func TestGenerateTokenZeroLength(t *testing.T) {
	tok := generateToken(0)
	if len(tok) != 32 {
		t.Errorf("expected default length 32, got %d", len(tok))
	}
}

func TestGenerateTokenURLSafe(t *testing.T) {
	for i := 0; i < 50; i++ {
		tok := generateToken(48)
		if strings.ContainsAny(tok, "+/") {
			t.Errorf("token contains non-URL-safe chars: %q", tok)
		}
	}
}

func TestRotateCmdRegistered(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "rotate <project> <VAR> [VAR...]" {
			return
		}
	}
	t.Error("rotate command not registered on rootCmd")
}
