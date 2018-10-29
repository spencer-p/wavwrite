package wavwrite

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

// Encode writes all audio streamed from s to w in WAVE format.
//
// Format precision must be 1 or 2 bytes. Format.DataSize must acccurately
// describe the bytes you will write.
func Encode(w io.Writer, s Streamer, format Format) (err error) {
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "wav")
		}
	}()

	if format.NumChannels <= 0 {
		return errors.New("invalid number of channels (less than 1)")
	}
	if format.Precision != 1 && format.Precision != 2 {
		return errors.New("unsupported precision, 1 or 2 is supported")
	}

	h := header{
		RiffMark:      [4]byte{'R', 'I', 'F', 'F'},
		FileSize:      int32(format.DataSize + 44), /* 44 is the header size */
		WaveMark:      [4]byte{'W', 'A', 'V', 'E'},
		FmtMark:       [4]byte{'f', 'm', 't', ' '},
		FormatSize:    16,
		FormatType:    1,
		NumChans:      int16(format.NumChannels),
		SampleRate:    int32(format.SampleRate),
		ByteRate:      int32(int(format.SampleRate) * format.NumChannels * format.Precision),
		BytesPerFrame: int16(format.NumChannels * format.Precision),
		BitsPerSample: int16(format.Precision) * 8,
		DataMark:      [4]byte{'d', 'a', 't', 'a'},
		DataSize:      int32(format.DataSize),
	}
	if err := binary.Write(w, binary.LittleEndian, &h); err != nil {
		return err
	}

	bw := bufio.NewWriter(w)
	written := 0
	for {
		n, err := s.Stream(bw)
		if err != nil {
			return err
		}
		if n == 0 {
			break
		}
		written += n
	}
	if err := bw.Flush(); err != nil {
		return err
	}

	// Check the bytes actually written matched what we said we would write
	if written != format.DataSize {
		return fmt.Errorf("format.DataSize (%d) does not match actual written bytes (%d)",
			format.DataSize, written)
	}

	return nil
}

type header struct {
	RiffMark      [4]byte
	FileSize      int32
	WaveMark      [4]byte
	FmtMark       [4]byte
	FormatSize    int32
	FormatType    int16
	NumChans      int16
	SampleRate    int32
	ByteRate      int32
	BytesPerFrame int16
	BitsPerSample int16
	DataMark      [4]byte
	DataSize      int32
}
