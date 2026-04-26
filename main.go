// Remove ANSI color escape sequences from the input stream.
//
// Usage:
//
//	cat FILE | discolor
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	BEL = 0x07
	LF  = 0x0A
	CR  = 0x0D
	ESC = 0x1B
)

func read(ch chan<- rune) {
	defer close(ch)

	r := bufio.NewReader(os.Stdin)

	for {
		c, _, err := r.ReadRune()

		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		ch <- c
	}
}

func main() {
	ch := make(chan rune)

	go read(ch)

	for c := range ch {
		switch c {
		case ESC:
			switch c = <-ch; c {
			case '[':
				for c := range ch {
					if c != ';' && c != '?' && (c < '0' || c > '9') {
						break
					}
				}
			case ']':
				if c = <-ch; c >= 0 && c <= '9' {
					for c := range ch {
						switch c {
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
			if c := <-ch; c != LF {
				_, _ = fmt.Printf("%c", c)
			}
		default:
			_, _ = fmt.Printf("%c", c)
		}
	}
}
