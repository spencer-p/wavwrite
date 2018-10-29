package wavwrite

// Streamer is able to stream a finite or infinite sequence of audio samples.
type Streamer interface {
	// Stream copies at most len(samples) next audio samples to the samples slice.
	//
	// The sample rate of the samples is unspecified in general, but should be specified for
	// each concrete Streamer.
	//
	// The value at samples[i][0] is the value of the left channel of the i-th sample.
	// Similarly, samples[i][1] is the value of the right channel of the i-th sample.
	//
	// Stream returns the number of streamed samples. If the Streamer is drained and no more
	// samples will be produced, it returns 0 and false. Stream must not touch any samples
	// outside samples[:n].
	//
	// There are 3 valid return pattterns of the Stream method:
	//
	//   1. n == len(samples) && ok
	//
	// Stream streamed all of the requested samples. Cases 1, 2 and 3 may occur in the following
	// calls.
	//
	//   2. 0 < n && n < len(samples) && ok
	//
	// Stream streamed n samples and drained the Streamer. Only case 3 may occur in the
	// following calls. If Err return a non-nil error, only this case is valid.
	//
	//   3. n == 0 && !ok
	//
	// The Streamer is drained and no more samples will come. Only this case may occur in the
	// following calls.
	Stream(samples []byte) (n int, ok bool)

	// Err returns an error which occured during streaming. If no error occured, nil is
	// returned.
	//
	// When an error occurs, Streamer must become drained and Stream must return 0, false
	// forever.
	//
	// The reason why Stream doesn't return an error is that it dramatically simplifies
	// programming with Streamer. It's not very important to catch the error right when it
	// happens.
	Err() error
}
