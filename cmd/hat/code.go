package main

import (
	"fmt"
	"io"
	"strings"
)

// readMasked forwards each line typed on keys to w, echoing it as ******** plus its last four
// characters, so a pasted sign-in code shows it arrived without being printed in full.
// keys is a terminal in raw mode; Ctrl+C calls cancel.
func readMasked(keys io.Reader, echo, w io.Writer, cancel func()) {
	var line []byte
	shown := 0
	buf := make([]byte, 256)
	for {
		n, err := keys.Read(buf)
		if err != nil {
			return
		}
		for _, b := range buf[:n] {
			switch {
			case b == '\r' || b == '\n':
				if len(line) == 0 {
					continue // second half of \r\n, or an empty Enter
				}
				fmt.Fprint(echo, "\r\n")
				w.Write(append(line, '\n'))
				line, shown = nil, 0
				continue
			case b == 3:
				cancel()
				return
			case b == 8 || b == 127:
				if len(line) > 0 {
					line = line[:len(line)-1]
				}
			case b >= ' ':
				line = append(line, b)
			}
			s := mask(line)
			fmt.Fprint(echo, strings.Repeat("\b \b", shown)+s)
			shown = len(s)
		}
	}
}

func mask(b []byte) string {
	if len(b) <= 4 {
		return strings.Repeat("*", len(b))
	}
	return "********" + string(b[len(b)-4:])
}
