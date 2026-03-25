// Remove ANSI color escape sequences from the input stream.
//
// Usage:
//
//	cat FILE | uncolor
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	BEL = 7
	LF  = 10
	CR  = 13
	ESC = 27
)

func main() {
	ch := make(chan rune)

	go func(ch chan<- rune) {
		br := bufio.NewReader(os.Stdin)

		for {
			r, _, err := br.ReadRune()

			if err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}

			ch <- r
		}

		close(ch)
	}(ch)

	for r := range ch {
		switch r {
		case ESC:
			switch r = <-ch; r {
			case '[':
				for r := range ch {
					if r != ';' && r != '?' && (r < '0' || r > '9') {
						break
					}
				}

			case ']':
				if r = <-ch; r >= 0 && r <= '9' {
					for r := range ch {
						switch r {

						case BEL:
							break

						case ESC:
							<-ch
							break
						}
					}
				}

			case '(', ')', '%':
				<-ch
			}

		case CR:
			if r := <-ch; r != LF {
				_, _ = fmt.Fprintf(os.Stdout, "%c", r)
			}

		default:
			_, _ = fmt.Fprintf(os.Stdout, "%c", r)
		}
	}
}
