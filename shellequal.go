package main

import (
	"bytes"
)

// routines used to figure out if the godownloader output script
// has changed or not

// stripShell takes a shell script and strips away
// all blank lines, all comments, and leading/trailing whitespace
//
func shellStrip(a []byte) []byte {
	lines := bytes.Split(a, []byte{'\n'})
	out := make([][]byte, 0, len(lines))
	for _, line := range lines {
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		out = append(out, bytes.TrimSpace(line))
	}
	return bytes.Join(out, []byte{'\n'})
}

// cmpShell returns true if two shell scripts are the same minus
// comments and whitespace.
//
// In future this might run mvdan/shfmt on the two scripts.
//
func shellEqual(a, b []byte) bool {
	return bytes.Equal(shellStrip(a), shellStrip(b))
}
