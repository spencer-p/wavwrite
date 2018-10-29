package wavwrite

import (
	"time"
)

// SampleRate is the number of samples per second.
type SampleRate int

// D returns the duration of n samples.
func (sr SampleRate) D(n int) time.Duration {
	return time.Second * time.Duration(n) / time.Duration(sr)
}

// N returns the number of samples that last for d duration.
func (sr SampleRate) N(d time.Duration) int {
	return int(d * time.Duration(sr) / time.Second)
}

// Format is the format of a Buffer or another audio source.
type Format struct {
	// SampleRate is the number of samples per second.
	SampleRate SampleRate

	// NumChannels is the number of channels. The value of 1 is mono, the value of 2 is stereo.
	// The samples should always be interleaved.
	NumChannels int

	// Precision is the number of bytes used to encode a single sample. Only values up to 6 work
	// well, higher values loose precision due to floating point numbers.
	Precision int
}

// Width returns the number of bytes per one frame (samples in all channels).
//
// This is equal to f.NumChannels * f.Precision.
func (f Format) Width() int {
	return f.NumChannels * f.Precision
}
