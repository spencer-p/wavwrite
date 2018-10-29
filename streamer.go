package wavwrite

import (
	"bufio"
)

// Streamer is able to stream a finite or infinite sequence of audio samples.
type Streamer interface {
	// Stream writes bytes to the described writer. It returns n, the amount of
	// bytes written.
	// Return n = 0 when there is nothing to write.
	// Return an error as necessary.
	Stream(w *bufio.Writer) (n int, err error)
}
