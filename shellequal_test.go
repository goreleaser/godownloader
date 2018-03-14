package main

import (
	"testing"
)

func TestShellStrip(t *testing.T) {
	var (
		src = `
#!/bin/sh

if [ "$1" = "echo" ]; then
    echo "echo"  
fi
`
		want = `if [ "$1" = "echo" ]; then
echo "echo"
fi`
	)
	got := string(shellStrip([]byte(src)))
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	if !shellEqual([]byte(src), []byte(want)) {
		t.Errorf("should be equal %q %q", src, want)
	}
}
