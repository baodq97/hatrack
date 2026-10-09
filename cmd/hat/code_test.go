package main

import (
	"strings"
	"testing"
)

func TestReadMasked(t *testing.T) {
	var echo, w strings.Builder
	readMasked(strings.NewReader("abc\x7fdefghij\r\nxy\r\n"), &echo, &w, func() { t.Fatal("cancelled") })
	if w.String() != "abdefghij\nxy\n" {
		t.Fatalf("forwarded %q", w.String())
	}
	if !strings.Contains(echo.String(), "********ghij\r\n") || strings.Contains(echo.String(), "abdef") {
		t.Fatalf("echo %q", echo.String())
	}

	cancelled := false
	w.Reset()
	readMasked(strings.NewReader("ab\x03cd\r"), &echo, &w, func() { cancelled = true })
	if !cancelled || w.Len() != 0 {
		t.Fatalf("ctrl+c: cancelled=%v forwarded %q", cancelled, w.String())
	}
}
